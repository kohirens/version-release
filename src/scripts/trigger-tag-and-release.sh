TriggerTagAndRelease() {
    hasTag=$(git show-ref "${CIRCLE_SHA1}" || echo "not found")
    # Skip if this commit is already tagged.
    if [ "${hasTag}" != "not found" ]; then
        echo "exiting, commit ${CIRCLE_SHA1} is already tagged: ${hasTag}"
        exit 1
    fi


    PARAM_MAP="{\"run-tag-and-release\": true, \"triggered-by-bot\": true, \"trigger-by\": \"${ORB_TRIGGER_NAME}}\" }"
    echo "{\"branch\": \"${PARAM_BRANCH}\", \"parameters\": ${PARAM_MAP}}" > pipelineparams.json
    cat pipelineparams.json
    TriggerPipeline
    Result
}

# Example API URL: https://circleci.com/api/v2/project/vcs-slug/org-name/repo-name/pipeline
# Note: keep in mind that you have to use a personal API token; project tokens are currently not supported on CircleCI API (v2) at this time.
# see: https://circleci.com/docs/api/v2/#operation/triggerPipeline
TriggerPipeline() {
    T=$(eval echo "$TOKEN")
    curl -u "${T}": -X POST --header "Content-Type: application/json" -d @pipelineparams.json \
      "${CIRCLECI_API_HOST}/api/v2/project/${PARAM_VCS_TYPE}/${CIRCLE_PROJECT_USERNAME}/${CIRCLE_PROJECT_REPONAME}/pipeline" -o /tmp/curl-result.txt
}

Result() {
    CURL_RESULT=$(cat /tmp/curl-result.txt)
    if [[ $(echo "$CURL_RESULT" | jq -r .message) == "Not Found" || $(echo "$CURL_RESULT" | jq -r .message) == "Permission denied" || $(echo "$CURL_RESULT" | jq -r .message) == "Project not found" ]]; then
        echo "Was unable to trigger tag-and-release workflow. API response: $(jq -r .message < /tmp/curl-result.txt)"
        exit 1
    else
        echo "Pipeline triggered!"
        echo "${CIRCLECI_APP_HOST}/jobs/${PARAM_VCS_TYPE}/${CIRCLE_PROJECT_USERNAME}/${CIRCLE_PROJECT_REPONAME}/$(jq -r .number < /tmp/curl-result.txt)"
    fi
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*$ORB_TEST_ENV}" == "$0" ]; then
    TriggerTagAndRelease
fi
