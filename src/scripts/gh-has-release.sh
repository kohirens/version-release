#!/bin/sh

set -e

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
        -w '%{http_code}' )

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
    has_release "${PARAM_OWNER}" "${PARAM_REPO}" "${PARAM_TAG}"
fi
