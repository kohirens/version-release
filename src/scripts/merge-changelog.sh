MergeChangelog() {
    changelogUpdated=$(git diff --name-only -- "${PARAM_CHANGELOG_FILE}")
    echo "BEGIN DEBUG"
    ls -la .
    git status
    echo "END DEBUG"

    if [ -z "${changelogUpdated}" ]; then
        echo "no changes detected in the ${PARAM_CHANGELOG_FILE} file"
        exit 0
    fi
    GEN_BRANCH_NAME="updated-changelog-skip-ci"
    git add CHANGELOG.md
    git config --global user.name "${CIRCLE_USERNAME}"
    git config --global user.email "${CIRCLE_USERNAME}@users.noreply.github.com"
    git checkout -b "${GEN_BRANCH_NAME}"
    git commit -m "Updated the ${PARAM_CHANGELOG_FILE}\n\n[skip ci]"
    # Do not run when sourced for bats-core
    if [ "${0#*$ORB_TEST_ENV}" == "$0" ]; then
        git push origin "${GEN_BRANCH_NAME}"
        echo "${GH_TOKE}" > really-i-need-a-file.txt
        gh auth login --with-token < really-i-need-a-file.txt
        gh pr create --base "${PARAM_BRANCH}" --head "${GEN_BRANCH_NAME}" --fill
        sleep 5
        gh pr merge --auto "--${PARAM_MERGE_TYPE}"
    fi
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*$ORB_TEST_ENV}" == "$0" ]; then
    MergeChangelog
fi
