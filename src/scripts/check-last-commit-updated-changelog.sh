# Check if the last commit contains the automated update message
CheckLastCommitUpdatedChangelog() {
    local isChangelogUpToDate
    local lastExitCode

    # check if the last commit message is the auto update message.
    isChangelogUpToDate=$(git show -s HEAD --grep 'automated update of CHANGELOG.md')
    lastExitCode=$?
    echo "lastExitCode=\"${lastExitCode}\""
    if [ -z "${isChangelogUpToDate}" ]; then
        echo "not ready for a release, exiting gracefully"
        exit 0
    fi
    echo "last commit indicates it was from auto updating the changelog"
    echo "last commit message:"
    echo "${isChangelogUpToDate}"
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*$ORB_TEST_ENV}" == "$0" ]; then
    CheckLastCommitUpdatedChangelog
fi
