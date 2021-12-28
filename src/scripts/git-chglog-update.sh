GitChglogUpdate() {
    hasTag=$(git tag)

    printf "hasTag = \"%s\"" "$hasTag"

    if [ -z "${hasTag}" ]; then
        git-chglog --output "${PARAM_OUTPUTFILE}" -c "${PARAM_CONFIGFILE}" --next-tag=0.1.0
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
