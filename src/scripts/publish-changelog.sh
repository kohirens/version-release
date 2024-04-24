publish_changelog() {
    vro publish-changelog \
        "${PARAM_CHANGELOG_FILE}" \
        "${PARAM_MAIN_TRUNK_BRANCH}" \
        "${PARAM_WORKING_DIRECTORY}"
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*"${ORB_TEST_ENV}"}" = "$0" ]; then
    publish_changelog
fi
