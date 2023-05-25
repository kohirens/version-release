# All functions that generate a workflow should:
# 1. Generate a PARAM_GENERATED_WORKFLOW_PATH
# 2. Generate a parameter JSON
# 3. Exit 0/1 with success/failure on generating the workflow files.

# Show debug information.
function debugGeneratedWorkflow() {
    ls -la .circleci
    echo
    echo "content of ${1}"
    echo "$(<"${1}")"
}

# Generate a Publish Changelog Workflow
generatePublishChangelogWorkflow() {
    local changelogUpdated
    local changelogUntracked

    changelogUpdated=$(git diff --name-only -- "${PARAM_CHANGELOG_PATH}")
    changelogUntracked=$(git status | grep "${PARAM_CHANGELOG_PATH}" || echo "")

    echo "changelogUpdated = \"${changelogUpdated}\""
    echo "changelogUntracked = \"${changelogUntracked}\""

    if [ -n "${changelogUntracked}" ]; then
        echo "Generating a publish changelog workflow"
        cp .circleci/publish-changelog.yml "${PARAM_GENERATED_WORKFLOW_PATH}"
        echo '{ "continued_automation": "publish-changelog" }' > "${PARAM_CONTINUE_PARAMETERS_PATH}"
        exit 0
    fi
}

# Generate Tag and Release Workflow
generateTagAndReleaseWorkflow() {
    local currVersion
    local revRange
    local isTaggable

    git-tool-belt semver -save build-version.json
    currVersion=$(jq -r .currentVersion < build-version.json)
    revRange="${currVersion}"

    if [ "${currVersion}" != "HEAD" ]; then
        revRange="${currVersion}..HEAD"
    fi

    isTaggable=$(git-tool-belt taggable --commitRange "${revRange}")

    echo "commit range ${revRange} tag ability is \"${isTaggable}\""

    if [ "${isTaggable}" = "true" ]; then
        echo "Generate tag and release workflow"
        cp .circleci/tag-and-release.yml "${PARAM_GENERATED_WORKFLOW_PATH}"
        echo '{ "continued_automation": "publish-new-tag" }' > "${PARAM_CONTINUE_PARAMETERS_PATH}"
        exit 0
    fi
}

# Generate a workflow for continuation.
GenerateAWorkflow() {
    # check if the changelog should be updated
    generatePublishChangelogWorkflow

    generateTagAndReleaseWorkflow
    if [ -n "${PARAM_TRIGGERED_FLOW}" ]; then
        echo "triggering workflow ${PARAM_TRIGGERED_FLOW}"
        echo "{ \"triggered_flow\": \"${PARAM_TRIGGERED_FLOW}\" }" > "${PARAM_CONTINUE_PARAMETERS_PATH}"
    fi
    if [ ! -f "${PARAM_GENERATED_WORKFLOW_PATH}" ] || [ ! -f "${PARAM_CONTINUE_PARAMETERS_PATH}" ]; then
        echo "Failed to generate the workflow files, cannot continue."
        exit 1
    fi
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*"${ORB_TEST_ENV}"}" = "$0" ]; then
    GenerateAWorkflow
    debugGeneratedWorkflow "${PARAM_CONTINUE_PARAMETERS_PATH}"
fi
