trigger_workflow() {
    # Get the value of the semantic version tag in 1 of 3 way.
    # If more than 1 is set, the last wins.
    semver=""
    if [ -n "${PARAM_TAG_CMD}" ]; then
        semver="$("${PARAM_TAG_CMD}")"
        echo "semantic version ${semver} was set by command"
    fi

    if [ -n "${PARAM_TAG_ENV_VAR}" ]; then
        semver="${!PARAM_TAG_ENV_VAR}"
        echo "semantic version ${semver} was extracted from environment variable ${PARAM_TAG_ENV_VAR}"
    fi

    if [ -n "${PARAM_TAG_FILE}" ]; then
        semver="$(cat "${PARAM_TAG_FILE}")"
        echo "semantic version ${semver} was pull from file ${PARAM_TAG_FILE}"
    fi

    if [ -n "${semver}" ]; then
        vro -semver "${semver}" workflow-selector \
            "${PARAM_CHANGELOG_FILE}" \
            "${PARAM_MAIN_TRUNK_BRANCH}" \
            "${PARAM_WORKING_DIRECTORY}" \
            "${CIRCLE_SHA1}"
    else
        vro workflow-selector \
            "${PARAM_CHANGELOG_FILE}" \
            "${PARAM_MAIN_TRUNK_BRANCH}" \
            "${PARAM_WORKING_DIRECTORY}" \
            "${CIRCLE_SHA1}"
    fi
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*"${ORB_TEST_ENV}"}" = "$0" ]; then
    trigger_workflow
fi
