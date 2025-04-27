#!/bin/bash

publish_tag_and_release() {
    # Get the value of the semantic version tag in 1 of 3 way.
    # If more than 1 is set, the last wins.
    semver=""

    if [ -n "${PARAM_TAG_CMD}" ]; then
        semver="$(eval "${PARAM_TAG_CMD}")"
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

    if [ "${CIRCLECI}" = "true" ]; then
        CICD_PLATFORM="circleci"
    fi

    if [ "${GITHUB_ACTIONS}" = "true" ]; then
        CICD_PLATFORM="github"
    fi

    cmd_str="avr"
    # Global options
    cmd_str="${cmd_str} -branch ""${PARAM_MAIN_TRUNK_BRANCH}"""
    cmd_str="${cmd_str} -cicd ""${CICD_PLATFORM}"""
    cmd_str="${cmd_str} -wd ""${PARAM_WORKING_DIRECTORY}"""
    cmd_str="${cmd_str} -github-api-url ""${PARAM_GITHUB_API_URL}"""
    cmd_str="${cmd_str} -github-server ""${PARAM_GITHUB_SERVER}"""

    if [ -n "${semver}" ]; then
        cmd_str="${cmd_str} -semver \"${semver}\""
    fi

    echo "PARAM_ENABLE_TAG_V_PREFIX=${PARAM_ENABLE_TAG_V_PREFIX}"
    if [ "${PARAM_ENABLE_TAG_V_PREFIX}" = "true" ] || [ "${PARAM_ENABLE_TAG_V_PREFIX}" = "1" ]; then
        cmd_str="${cmd_str} -enable-tag-v-prefix"
    fi

    # sub command
    cmd_str="${cmd_str} publish-release-tag"
    # arguments
    cmd_str="${cmd_str} ""${PARAM_MAIN_TRUNK_BRANCH}"""

    sh -c "${cmd_str}"
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*"${ORB_TEST_ENV}"}" = "$0" ]; then
    publish_tag_and_release
fi
