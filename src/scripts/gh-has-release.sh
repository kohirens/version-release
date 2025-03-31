#!/bin/bash

has_release() {
    owner_slash_repo="${1}"
    semver_tag="${2}"
    gh_token="${3}"

    http_code=$(curl -kL \
        --show-error \
        --silent \
        -H "Accept: application/vnd.github+json" \
        -H "Authorization: Bearer ${gh_token}" \
        -H "X-GitHub-Api-Version: 2022-11-28" \
        --url "https://api.github.com/repos/${owner_slash_repo}/releases/tags/${semver_tag}" \
        --output /dev/null \
        -w '%{http_code}')

    echo "${http_code}"
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*"${ORB_TEST_ENV}"}" = "$0" ]; then
    if [ "${GITHUB_ACTIONS}" != "true" ]; then
        echo "only gitHub actions should run this script, abort!"
        exit 1
    fi

    semver=""
    if [ -n "${PARAM_TAG_CMD}" ]; then
        semver="$(eval "${PARAM_TAG_CMD}")"
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

    result="$(has_release "${PARAM_OWNER_SLASH_REPO}" "${semver}" "${GH_WRITE_TOKEN}")"

    echo "has_release result: \"${result}\""

    if [ "${result}" = "200" ]; then
        printf "aborting... there is already a release tag at https://github.com/%s/releases/tag/%s\n" "${PARAM_OWNER_SLASH_REPO}" "${semver}"
        exit 1
    else
        printf "will proceed since no exiting release tag was found at https://github.com/%s/releases/tag/%s\n" "${PARAM_OWNER_SLASH_REPO}" "${semver}"
    fi
fi
