TagAndRelease() {
    git fetch --all -p

    if [ -f "${PARAM_COMMIT_FILE}" ]; then
        currCommit=$(cat "${PARAM_COMMIT_FILE}")
        if [ "${currCommit}" != "" ]; then
          echo "resetting to commit ${currCommit}"
          git reset --hard "${currCommit}"
        fi
    fi

    changelogUpdated=$(git diff --name-only -- "${PARAM_CHANGELOG_FILE}")
    echo "changelogUpdated=${changelogUpdated}"
    # Skip if there are changes in the changelog that have not been merged.
    if [ -n "${changelogUpdated}" ]; then
        echo "exiting, changes detected in the ${PARAM_CHANGELOG_FILE} file"
        # ensure this file exist so persist does not fail.
        echo "" > "${PARAM_TAG_FILE}"
        exit 0
    fi

    hasTag=$(git show-ref "${CIRCLE_SHA1}" || echo "not found")
    echo "hasTag=${hasTag}"
    # Skip if this commit is already tagged.
    if [ "${hasTag}" != "not found" ]; then
        echo "exiting, commit is already tagged: ${hasTag}"
        # ensure this file exist so persist does not fail.
        echo "" > "${PARAM_TAG_FILE}"
        exit 0
    fi

    git-tool-belt semver -save build-version.json
    nextVersion=$(jq -r .nextVersion < build-version.json)
    currVersion=$(jq -r .currentVersion < build-version.json)
    releaseDay=$(date +"%Y-%m-%d")
    revRange="${currVersion}"
    # currVersion defaults to "HEAD" if it cannot find a tag, grabbing and tagging the whole history up to this point
    if [ "${currVersion}" != "HEAD" ]; then
        revRange="${currVersion}..HEAD"
    fi

    # Fetch all the refs
    isTaggable=$(git-tool-belt taggable --commitRange "${revRange}")
    echo "commit range ${revRange} tag ability is \"${isTaggable}\""
    # Skip if there are no notable commits to tag.
    if [ "${isTaggable}" = "false" ]; then
        echo "exiting, no notable commits to tag"
        # ensure this file exist so persist does not fail.
        echo "" > "${PARAM_TAG_FILE}"
        exit 0
    fi

    echo "tagging commit hash ${CIRCLE_SHA1} with ${nextVersion}"
    # Switch to SSH to use the token stored in the environment.
    gh config set git_protocol ssh --host github.com
    gh auth status --hostname github.com
    gh release create "${nextVersion}" --generate-notes --target "${PARAM_BRANCH}" --title "${nextVersion} - ${releaseDay}"
    echo "${nextVersion}" > "${PARAM_TAG_FILE}"
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*$ORB_TEST_ENV}" == "$0" ]; then
    TagAndRelease
fi
