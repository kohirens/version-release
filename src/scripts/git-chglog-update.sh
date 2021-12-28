GitChglogUpdate() {
    wd=$(pwd)
    printf "working directory: %s\n" "${wd}"
    git-chglog --output "${PARAM_OUTPUTFILE}" -c "${PARAM_CONFIGFILE}"
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*$ORB_TEST_ENV}" == "$0" ]; then
    GitChglogUpdate
fi
