trigger_workflow() {
    vro workflow-selector \
        "${PARAM_CHANGELOG_FILE}" \
        "${PARAM_MAIN_TRUNK_BRANCH}" \
        "${PARAM_WORKING_DIRECTORY}" \
        "${CIRCLE_SHA1}"
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*"${ORB_TEST_ENV}"}" = "$0" ]; then
    trigger_workflow
fi
