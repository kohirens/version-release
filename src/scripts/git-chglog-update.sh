GitChglogUpdate() {
    hasTag=$(git tag)

    if [ -z "${hasTag}" ]; then
        echo "no tags found"
        git-chglog --output "${PARAM_OUTPUTFILE}" -c "${PARAM_CONFIGFILE}" --next-tag=0.1.0
    elif [ -f "/usr/local/bin/git-tool-belt" ]; then
        git-tool-belt semver -save build-version.json
        nextVersion=$(jq -r .nextVersion < build-version.json)
        echo "nextVersion = ${nextVersion}"
        git-chglog --output "${PARAM_OUTPUTFILE}" -c "${PARAM_CONFIGFILE}" --next-tag="${nextVersion}"
    else
        wd=$(pwd)
        printf "updating %s using config file %s\n" "${wd}/${PARAM_OUTPUTFILE}" "${PARAM_CONFIGFILE}"
        git-chglog --output "${PARAM_OUTPUTFILE}" -c "${PARAM_CONFIGFILE}"
    fi
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*$ORB_TEST_ENV}" == "$0" ]; then
    GitChglogUpdate
fi
