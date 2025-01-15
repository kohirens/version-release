# Runs prior to every test
setup() {
    source ./src/scripts/gh-has-release.sh
}

@test '1: has release' {
    # setup
    export GITHUB_ACTIONS="true"
    export PARAM_OWNER_SLASH_REPO="kohirens/version-release"
    export PARAM_TAG="0.1.0"

    # test
    result=$(has_release "${PARAM_OWNER_SLASH_REPO}" "${PARAM_TAG}")

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
    result=$(has_release "${PARAM_OWNER_SLASH_REPO}" "${PARAM_TAG}")

    # assert
    [ "${result}" = "404" ]
}
