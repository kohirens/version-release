publish_tag_and_release() {
    vro publish-release-tag \
        "${PARAM_CHANGELOG_FILE}" \
        "${PARAM_MAIN_TRUNK_BRANCH}" \
        "${PARAM_WORKING_DIRECTORY}"
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*"${ORB_TEST_ENV}"}" = "$0" ]; then
    publish_tag_and_release
fi
