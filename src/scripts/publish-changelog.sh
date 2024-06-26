publish_changelog() {
    # publishing a changelog branch uses git push, so we will run into
    # known host issues with github.com servers from time to time, this adds
    # the latest github.com public SSH keys to the known_hosts files to resolve
    # that issue.
    if [ -f ~/.ssh/known_hosts ]; then
        vro known-sshkeys >> ~/.ssh/known_hosts
    else
        mkdir -p ~/.ssh
        vro known-sshkeys > ~/.ssh/known_hosts
    fi

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
        echo "semantic version ${semver} was pulled from the file ${PARAM_TAG_FILE}"
    fi

    if [ -n "${semver}" ]; then
        vro -semver "${semver}" \
            publish-changelog \
            "${PARAM_CHANGELOG_FILE}" \
            "${PARAM_MAIN_TRUNK_BRANCH}" \
            "${PARAM_WORKING_DIRECTORY}"
    else
        vro publish-changelog \
            "${PARAM_CHANGELOG_FILE}" \
            "${PARAM_MAIN_TRUNK_BRANCH}" \
            "${PARAM_WORKING_DIRECTORY}"
    fi
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*"${ORB_TEST_ENV}"}" = "$0" ]; then
    publish_changelog
fi
