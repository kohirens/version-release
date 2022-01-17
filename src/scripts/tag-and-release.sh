TagAndRelease() {
    changelogUpdated=$(git diff --name-only -- "${PARAM_CHANGELOG_FILE}")
    # Skip if there are changes in the changelog that have not been merged.
    if [ -n "${changelogUpdated}" ]; then
        echo "exiting, changes detected in the ${PARAM_CHANGELOG_FILE} file"
        exit 0
    fi

    hasTag=$(git show-ref --tags | grep "${CIRCLE_SHA1}" | grep "refs/tags/v\\?\\d\\.\\d\\.\\d")
    # Skip if this commit is already tagged.
    if [ -n "${hasTag}" ]; then
        echo "exiting, commit is already tagged: ${hasTag}"
        exit 0
    fi

    git-tool-belt version
    nextVersion=$(jq -r .nextVersion < build-version.json)
    prevVersion=$(jq -r .currentVersion < build-version.json)
    releaseDay=$(date +"%Y-%m-%d")

    isTaggable=$(git-tool-belt taggable "${prevVersion}..HEAD")
    echo "${prevVersion}..HEAD isTaggable = ${isTaggable}"
    # Skip if there are no notable commits to tag.
    if [ "${isTaggable}" = "true" ]; then
        echo "exiting, no notable commits to tag"
        exit 0
    fi

    echo "tagging commit hash ${CIRCLE_SHA1} with ${nextVersion}"
    echo "${GH_TOKE}" > really-i-need-a-file.txt
    gh auth login --with-token < really-i-need-a-file.txt
    gh release create "${nextVersion}" --generate-notes --target "${PARAM_BRANCH}" --title "[${nextVersion}] - ${releaseDay}"

}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*$ORB_TEST_ENV}" == "$0" ]; then
    TagAndRelease
fi
