MergeChangelog() {
    prevCommit=$(git rev-parse origin/"${PARAM_BRANCH}")
    # make the file so that the persist step does not fail.
    echo "" > "${PARAM_COMMIT_FILE}"

    changelogUpdated=$(git diff --name-only -- "${PARAM_CHANGELOG_FILE}")

    if [ -z "${changelogUpdated}" ]; then
        echo "no changes detected in the ${PARAM_CHANGELOG_FILE} file"
        exit 0
    fi
    GEN_BRANCH_NAME="updated-changelog-skip-ci"
    git add CHANGELOG.md
    git config --global user.name "${CIRCLE_USERNAME}"
    git config --global user.email "${CIRCLE_USERNAME}@users.noreply.github.com"
    git checkout -b "${GEN_BRANCH_NAME}"
    gentBranchCommitMsg="Updated the ${PARAM_CHANGELOG_FILE} [skip ci]"
    git commit -m "${gentBranchCommitMsg}"
    # Do not run when sourced for bats-core
    if [ "${0#*$ORB_TEST_ENV}" == "$0" ]; then
        git push origin "${GEN_BRANCH_NAME}"
        echo "${GH_TOKE}" > really-i-need-a-file.txt
        gh auth login --with-token < really-i-need-a-file.txt
        gh pr create --base "${PARAM_BRANCH}" --head "${GEN_BRANCH_NAME}" --fill
        sleep 5
        gh pr merge --auto "--${PARAM_MERGE_TYPE}"

        waitForPrToMerge
    fi
}

waitForPrToMerge() {
    # Wait until the branch is fully merged. and the merge branch has been updated.
    # This will help make sure that operations started in this job complete before moving on.
    # 1. Record the commit has from committing the changelog
    prCommit=$(git rev-parse "${GEN_BRANCH_NAME}")
    # 2. Fetch remote changes
    git fetch --all -p
    echo "prevCommit = ${prevCommit}"
    echo "prCommit = ${prCommit}"
    printf "merging pr is "
    # 3. Loop for so many seconds
    counter=0
    while [ $counter -lt 10 ]; do
        # 4. Get the latest commit of the merge branch
        currCommit=$(git rev-parse origin/"${PARAM_BRANCH}")
        echo "currCommit = ${currCommit}"
        # 5. Check to see if the merge branch previous commit and current commit have changed.
        if [ "${currCommit}" != "${prevCommit}" ]; then
            echo " done"
            currCommitMsg=$(git show-branch --no-name "${currCommit}")
            if [ "${currCommitMsg}" = "${gentBranchCommitMsg}" ]; then
                echo "merge has completed successfully"
                echo "${currCommit}" > "${PARAM_COMMIT_FILE}"
                cm=$(cat "${PARAM_COMMIT_FILE}")
                printf "cm = %s"$'\n' "${cm}"
            fi
        else
            printf "."
        fi
        counter=$((counter+1))
        # wait a second
        sleep 1
    done
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*$ORB_TEST_ENV}" == "$0" ]; then
    MergeChangelog
fi
