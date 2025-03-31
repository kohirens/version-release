
# Runs prior to every test
setup() {
    source ./src/tests/bin/install.sh
    InstallPseudoCmd "curl"
    source ./src/scripts/gh-has-release.sh
}

@test '1: has release' {
    # setup
    export GITHUB_ACTIONS="true"
    export PARAM_OWNER_SLASH_REPO="kohirens/version-release"
    export PARAM_TAG="4.2.4"

    # test
    result=$(has_release "${PARAM_OWNER_SLASH_REPO}" "${PARAM_TAG}" "fake-token")

    echo "has_release result: \"${result}\""

    # assert
    [ "${result}" = "200" ]
}

@test '2: does not have such a release' {
    # setup
    export GITHUB_ACTIONS="true"
    export PARAM_OWNER_SLASH_REPO="kohirens/version-release"
    export PARAM_TAG="f4.2.4-fake"

    # test
    result=$(has_release "${PARAM_OWNER_SLASH_REPO}" "${PARAM_TAG}" "fake-token")

    # assert
    [ "${result}" = "404" ]
}
