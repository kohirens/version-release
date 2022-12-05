TagAndRelease() {
    # pre-checks
    if [ -z "${GH_TOKEN}${PARAM_GH_TOKEN_VAR}" ]; then
        echo "The environment variable name that should point to a GitHub write token is empty."
        echo "Please set the docs for the tag-and-release job parameter \"gh_token_var\" and try again."
        exit 1
    fi

    if [ -n "${PARAM_GH_TOKEN_VAR}" ]; then
        # seems to be the best way for connecting to the Github using the CLI.
        export GH_TOKEN="${!PARAM_GH_TOKEN_VAR}"
    fi

    # require a GH_TOKEN
    if [ -z "${GH_TOKEN}" ]; then
        echo "No GitHub write token found."
        echo "Please set the environment variable GH_TOKEN."
        echo "You can also specify which environment variable to use, see \"gh_token_var\" parameter in the tag-and-release job."
        exit 1
    fi

    git fetch --all -p

    if [ -f "${PARAM_COMMIT_FILE}" ]; then
        currCommit=$(<"${PARAM_COMMIT_FILE}")
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

    # Switch to SSH to use the token stored in the environment variable GH_TOKEN.
    gh config set git_protocol ssh --host "${PARAM_GH_SERVER}"
    echo
    echo
    # see: https://josh-ops.com/posts/gh-auth-login-in-actions/
    # NOTE: Using just GH_TOKEN set in the environment seems to fail when you have
    #       to supply the --hostname flag to point to a GHE server.
    if [ "${PARAM_GH_SERVER}" != "github.com" ]; then
        echo "login to ${PARAM_GH_SERVER}"
        echo "${GH_TOKEN}" | gh auth login --hostname "${PARAM_GH_SERVER}" --with-token
        echo
        echo
    fi
    echo "auth status of ${PARAM_GH_SERVER}"
    gh auth status --hostname "${PARAM_GH_SERVER}"
    echo
    echo
    echo "tagging commit hash ${CIRCLE_SHA1} with ${nextVersion}"
    gh release create "${nextVersion}" --generate-notes --target "${PARAM_BRANCH}" --title "${nextVersion} - ${releaseDay}"
    echo "${nextVersion}" > "${PARAM_TAG_FILE}"
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*"${ORB_TEST_ENV}"}" = "$0" ]; then
    TagAndRelease
fi
