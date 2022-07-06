AddMissingChgLogConfig() {
    git-tool-belt checkconf -path "${PARAM_CONFIG_FILE}" -repo "${CIRCLE_REPOSITORY_URL}"
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*$ORB_TEST_ENV}" == "$0" ]; then
    AddMissingChgLogConfig
fi
