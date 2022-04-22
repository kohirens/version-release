MergeChangelog() {
    prevCommit=$(git rev-parse origin/"${PARAM_BRANCH}")

    # make the file so that the persist step does not fail.
    echo "" > "${PARAM_COMMIT_FILE}"

    changelogUpdated=$(git diff --name-only -- "${PARAM_CHANGELOG_FILE}")
    changelogUntracked=$(git status | grep "${PARAM_CHANGELOG_FILE}" || echo "")
    if [ "${changelogUntracked}" != "" ]; then
        # TODO: Test with a repo with tags but no CHANGELOG.md
        echo "added new ${PARAM_CHANGELOG_FILE} file"
    elif [ -z "${changelogUpdated}" ]; then
        echo "no changes detected in the ${PARAM_CHANGELOG_FILE} file"
        exit 0
    fi

    # Exit if branch exist remotely already.
    branchExistRemotely=$(git ls-remote --heads "${CIRCLE_REPOSITORY_URL}" "${GEN_BRANCH_NAME}" | wc -l)
    echo "branchExistRemotely = ${branchExistRemotely}"
    # Exit if branch exist remotely already.
    if [ "${branchExistRemotely}" = "1"  ]; then
        echo "the branch '${GEN_BRANCH_NAME}' exists on ${CIRCLE_REPOSITORY_URL}, please remove it manually so this job can complete successfully; exiting with code 1"
        exit 1
    fi

    GEN_BRANCH_NAME="updated-changelog-skip-ci"
    git add CHANGELOG.md
    git config --global user.name "${CIRCLE_USERNAME}"
    git config --global user.email "${CIRCLE_USERNAME}@users.noreply.github.com"
    git checkout -b "${GEN_BRANCH_NAME}"
    mergeBranchCommitMsg="Updated the ${PARAM_CHANGELOG_FILE}"
    git commit -m "${mergeBranchCommitMsg}" -m "[skip ci]"
    # Do not run when sourced for bats-core
    # TODO: This can be tested if you mock the gh, git, and setup a dummy repo at test time.
    if [ "${0#*$ORB_TEST_ENV}" == "$0" ]; then
        git push origin "${GEN_BRANCH_NAME}"
        # Switch to SSH to use the token stored in the environment.
        gh config set git_protocol ssh --host github.com
        gh auth status --hostname github.com
        gh pr create --base "${PARAM_BRANCH}" --head "${GEN_BRANCH_NAME}" --fill
        sleep 5
        gh pr merge --auto "--${PARAM_MERGE_TYPE}"

        waitForPrToMerge
    fi
}

waitForPrToMerge() {
    # Wait until the branch is fully merged. and the merge branch has been updated.
    # This will help make sure that operations started in this job complete before moving on.
    printf "%s" "merging pr is "
    # 1. Loop for so many seconds
    counter=0
    while [ $counter -lt 10 ]; do
        # 2. Fetch remote changes
        git fetch --all -p
        # 3. Get the latest commit of the merge branch
        currCommit=$(git rev-parse origin/"${PARAM_BRANCH}")
        # 4. Check to see if the merge branch previous commit and current commit have changed.
        if [ "${currCommit}" != "${prevCommit}" ]; then
            echo " done"
            currCommitMsg=$(git show-branch --no-name "${currCommit}")
            # 5. If the commit messages are the same, then make a file to persist to the next job
            if [ "${currCommitMsg}" = "${mergeBranchCommitMsg}" ]; then
                echo "merge has completed successfully"
                echo "${currCommit}" > "${PARAM_COMMIT_FILE}"
                cm=$(cat "${PARAM_COMMIT_FILE}")
                printf "cm = %s"$'\n' "${cm}"
            fi
            # Exit the loop
            break
        else
            printf "."
        fi
        counter=$((counter+1))
        # Wait a second
        sleep 1
    done
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*$ORB_TEST_ENV}" == "$0" ]; then
    MergeChangelog
fi
