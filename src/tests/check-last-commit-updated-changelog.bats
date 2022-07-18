# Runs prior to every test
setup() {
    # Install required software
    source ./src/scripts/installs.sh
    InstallGitChglog
    InstallGitToolBelt

    # Load our script file.
    source ./src/scripts/check-last-commit-updated-changelog.sh
}

UnbundleRepo() {
    git clone -b main src/tests/data/${1}.bundle "${2}"
}

@test '1: Test last commit contains automated changelog message' {
    local fixture
    local ret_code
    local subShell

    rm -rf "tmp/${fixture}"
    fixture="repo-04"
    # Extract repo
    UnbundleRepo "${fixture}" "tmp/${fixture}"
    cd "tmp/${fixture}"

    # Run the function in a way that allows us to capture the exit code.
    subShell=$(CheckLastCommitUpdatedChangelog)
    ret_code=$?
    # Test output
    echo
    echo "ret_code=${ret_code}"
    echo "${subShell}" | grep 'last commit indicates it was from auto updating the changelog'
}

@test '2: Test last commit does not contain automated changelog message' {
    local fixture
    local ret_code
    local subShell

    rm -rf "tmp/${fixture}"
    fixture="repo-05"
    # Extract repo
    UnbundleRepo "${fixture}" "tmp/${fixture}"
    cd "tmp/${fixture}"

    # Run the function in a way that allows us to capture the exit code.
    subShell=$(CheckLastCommitUpdatedChangelog)
    ret_code=$?
    # Test output
    echo
    echo "ret_code=${ret_code}"
    echo "${subShell}" | grep 'not ready for a release, exiting gracefully'
}
