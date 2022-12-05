TriggerTagAndRelease() {
    set -e

    # pre-checks
    if [ -z "${PARAM_GH_TOKEN}" ]; then
        echo "An environment variable name that should point to a GitHub write token is empty."
        echo "Please set the docs for the tag-and-release job parameter \"gh_token\" and try again."
        exit 1
    else
        # the seems to be the best way for connecting to the Github using the CLI.
        export GH_TOKEN="${!PARAM_GH_TOKEN}"
    fi

    if [ -z "${PARAM_CIRCLE_TOKEN}" ]; then
        echo "environment variable name that should point to a CircleCI token is empty."
        echo "Please set the docs for the tag-and-release job parameter \"circle_token_var\" and try again."
        exit 1
    fi

    if [ -z "${GH_TOKEN}" ]; then
        echo "No GitHub write token found."
        echo "Please set the environment variable GH_TOKEN."
        echo "You can also specify which environment variable to use, see \"gh_token_var\" parameter in the publish-changelog job."
        exit 1
    fi

    echo "checking for commits containing that conventional commit syntax..."

    # Do not trigger a release if there is nothing to tag.
    if [ ! -f "trigger.txt" ]; then
        exit 0
    fi
    # Again, exit gracefully if this was not meant to be triggered.
    TRIGGER_TXT=$(cat trigger.txt)
    if [ "${TRIGGER_TXT}" != "trigger-tag-and-release" ]; then
        exit 0
    fi

    hasTag=$(git show-ref "${CIRCLE_SHA1}" || echo "not found")
    # Skip if this commit is already tagged.
    if [ "${hasTag}" != "not found" ]; then
        echo "exiting, commit ${CIRCLE_SHA1} is already tagged: ${hasTag}"
        exit 1
    fi

    PARAM_MAP="{\"triggered_by_bot\": true}"
    echo "{\"branch\": \"${PARAM_BRANCH}\", \"parameters\": ${PARAM_MAP}}" > pipelineparams.json
    cat pipelineparams.json

    # BEGIN these work together
    if [ -z "${CIRCLECI_API_HOST}" ]; then
        echo "CIRCLECI_API_HOST environment variable is not set, please pass in 'circleci_api_host' parameter"
    fi
    if [ -z "${CIRCLECI_APP_HOST}" ]; then
        echo "CIRCLECI_APP_HOST environment variable is not set, please pass in 'circleci_app_host' parameter"
    fi
    if [ -z "${CIRCLECI_API_HOST}" ] || [ -z "${CIRCLECI_APP_HOST}" ]; then
        exit 1
    fi
    # END these work together

    TriggerPipeline
}

# Example API URL: https://circleci.com/api/v2/project/vcs-slug/org-name/repo-name/pipeline
# Note: keep in mind that you have to use a personal API token; project tokens are currently not supported on CircleCI API (v2) at this time.
# see: https://circleci.com/docs/api/v2/#operation/triggerPipeline
TriggerPipeline() {
    projectUrl="${CIRCLECI_API_HOST}/api/v2/project/${PARAM_VCS_TYPE}/${CIRCLE_PROJECT_USERNAME}/${CIRCLE_PROJECT_REPONAME}/pipeline"
    resultFile="/tmp/curl-result.txt"
    echo "circleci project URL: ${projectUrl}"

    curl -u "${!PARAM_CIRCLE_TOKEN}": -X POST --header "Content-Type: application/json" -d @pipelineparams.json \
      "${projectUrl}" -o "${resultFile}"

    resultMessage=$(jq -r .message < ${resultFile})

    if [ "${resultMessage}" = "Not Found" ] || [ "${resultMessage}" = "Permission denied" ] || [ "${resultMessage}" = "Project not found" ]; then
        echo "Was unable to trigger tag-and-release workflow. API response: $(jq -r .message < /tmp/curl-result.txt)"
        exit 1
    else
        echo "Pipeline triggered!"
        echo "${CIRCLECI_APP_HOST}/jobs/${PARAM_VCS_TYPE}/${CIRCLE_PROJECT_USERNAME}/${CIRCLE_PROJECT_REPONAME}/$(jq -r .number < ${resultFile})"
    fi
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*"${ORB_TEST_ENV}"}" = "$0" ]; then
    TriggerTagAndRelease
fi
