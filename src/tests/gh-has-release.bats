# Runs prior to every test
setup() {
    source ./src/scripts/gh-has-release.sh
}

@test '1: has release' {
    # setup
    export PARAM_OWNER="kohirens"
    export PARAM_REPO="version-release"
    export PARAM_TAG="4.2.4"

    # test
    result=$(has_release "${PARAM_OWNER}" "${PARAM_REPO}" "${PARAM_TAG}")

    # assert
    [ "${result}" = "yes" ]
}

@test '2: does not have such a release' {
    # setup
    export PARAM_OWNER="kohirens"
    export PARAM_REPO="version-release"
    export PARAM_TAG="f4.2.4-fake"

    # test
    result=$(has_release "${PARAM_OWNER}" "${PARAM_REPO}" "${PARAM_TAG}")

    # assert
    [ "${result}" = "no" ]
}
