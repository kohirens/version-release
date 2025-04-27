trigger_workflow() {
    # Get the value of the semantic version tag in 1 of 3 ways.
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
        echo "semantic version ${semver} was pulled from the file ${PARAM_TAG_FILE}"
    fi

    if [ -n "${CIRCLE_SHA1}" ]; then
        GIT_SHA="${CIRCLE_SHA1}"
    elif [ -n "${GITHUB_SHA}" ]; then
        GIT_SHA="${GITHUB_SHA}"
    fi

    if [ "${CIRCLECI}" = "true" ]; then
        CICD_PLATFORM="circleci"
    fi

    if [ "${GITHUB_ACTIONS}" = "true" ]; then
        CICD_PLATFORM="github"
    fi

    if [ -n "${semver}" ]; then
        avr \
            -branch "${PARAM_MAIN_TRUNK_BRANCH}" \
            -cicd "${CICD_PLATFORM}" \
            -semver "${semver}" \
            -wd "${PARAM_WORKING_DIRECTORY}" \
            workflow-selector \
                "${PARAM_CHANGELOG_FILE}" \
                "${GIT_SHA}"
    else
        avr \
            -branch "${PARAM_MAIN_TRUNK_BRANCH}" \
            -cicd "${CICD_PLATFORM}" \
            -wd "${PARAM_WORKING_DIRECTORY}" \
            workflow-selector \
                "${PARAM_CHANGELOG_FILE}" \
                "${GIT_SHA}"
    fi
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*"${ORB_TEST_ENV}"}" = "$0" ]; then
    trigger_workflow
fi
