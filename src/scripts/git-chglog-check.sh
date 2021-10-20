returnVal=0
GitChglogConfigCheck() {
    if [ ! -f "${PARAM_CONFIGFILE}" ]; then
        echo "no changelog file found."
        returnVal=1
    fi
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*$ORB_TEST_ENV}" == "$0" ]; then
    GitChglogConfigCheck
    exit $returnVal
fi
