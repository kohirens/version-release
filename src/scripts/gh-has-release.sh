has_release() {
    OWNER="${1}"
    REPO="${2}"
    TAG="${3}"

    http_code=$(curl -kL \
        --show-error \
        --silent \
        -H "Accept: application/vnd.github+json" \
        -H "Authorization: Bearer ${GH_TOKEN}" \
        -H "X-GitHub-Api-Version: 2022-11-28" \
        --url "https://api.github.com/repos/${OWNER}/${REPO}/releases/tags/${TAG}" \
        --output /dev/null \
        -w '%{http_code}')

    if [ "${http_code}" = "200" ]; then
        echo "yes"
    else
        echo "no"
    fi
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*"${ORB_TEST_ENV}"}" = "$0" ]; then
    semver=""
    if [ -n "${PARAM_TAG_CMD}" ]; then
        semver="$("${PARAM_TAG_CMD}")"
        echo "release tag ${semver} was set by command"
    fi

    if [ -n "${PARAM_TAG_ENV_VAR}" ]; then
        semver="${!PARAM_TAG_ENV_VAR}"
        echo "release tag ${semver} was extracted from environment variable ${PARAM_TAG_ENV_VAR}"
    fi

    if [ -n "${PARAM_TAG_FILE}" ]; then
        semver="$(cat "${PARAM_TAG_FILE}")"
        echo "release tag ${semver} was pulled from the file ${PARAM_TAG_FILE}"
    fi

    if [ -z "${semver}" ]; then
        echo "please set either parameter tag_cmd, tag_env_var, or tag_file to set the release tag to lookup"
        exit 1
    fi

    has_release "${PARAM_OWNER}" "${PARAM_REPO}" "${semver}" > "${PARAM_FILE}"

    result=$(cat "${PARAM_FILE}")

    if [ "${result}" = "yes" ]; then
        echo "release tag ${result} was found"
    else
        echo "cannot find a release tag ${result} at github.com/${PARAM_OWNER}/${PARAM_REPO}"
    fi
fi
