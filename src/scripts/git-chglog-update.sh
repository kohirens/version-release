GitChglogUpdate() {
    hasTag=$(git tag)

    cd "${PARAM_REPOSITORY_PATH}" || exit

    echo
    echo
    pwd
    echo
    echo

    if [ -z "${hasTag}" ]; then
        echo "no tags found"
        git-chglog --output "${PARAM_CHANGELOG_FILE}" -c "${PARAM_CONFIG_FILE}" --next-tag=0.1.0
    elif [ -f "/usr/local/bin/git-tool-belt" ]; then
        echo "git-tool-belt found"
        git-tool-belt semver -save build-version.json
        nextVersion=$(jq -r .nextVersion < build-version.json)
        echo "nextVersion = ${nextVersion}"
        git-chglog --output "${PARAM_CHANGELOG_FILE}" -c "${PARAM_CONFIG_FILE}" --next-tag="${nextVersion}"
    else
        echo "running git-chglog as normal"
        wd=$(pwd)
        printf "updating %s using config file %s\n" "${wd}/${PARAM_CHANGELOG_FILE}" "${PARAM_CONFIG_FILE}"
        git-chglog --output "${PARAM_CHANGELOG_FILE}" -c "${PARAM_CONFIG_FILE}"
    fi
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*"${ORB_TEST_ENV}"}" == "$0" ]; then
    GitChglogUpdate
fi
