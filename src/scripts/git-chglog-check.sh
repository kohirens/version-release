GitChglogConfigCheck() {
    if [ -f "${PARAM_CONFIGFILE}" ]; then
        return 0
    fi

    return 1
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*$ORB_TEST_ENV}" == "$0" ]; then
    GitChglogConfigCheck
    exit $?
fi
