AddMissingChgLogConfig() {
    wd=$(pwd)
    echo "working directory = ${wd}"

    if [ ! -f "${PARAM_CONFIGFILE}" ]; then
        CONFIG_DIR=$(dirname "${PARAM_CONFIGFILE}")
        echo "no changelog file found at ${CONFIG_DIR}. adding the default..."
        mkdir -p "${CONFIG_DIR}"
        sed -e "s/\${VCS_URL}/${CIRCLE_REPOSITORY_URL}/" ./files/config.yml > "${PARAM_CONFIGFILE}"
        cp ./files/CHANGELOG.tpl.md "${CONFIG_DIR}"

        if [ ! -f "${PARAM_CONFIGFILE}" ]; then
            echo "unable to add ${PARAM_CONFIGFILE}"
        fi
    else
        echo "found git-chglog config here ${PARAM_CONFIGFILE}"
    fi
}

# Will not run if sourced for bats-core tests.
# View src/tests for more information.
ORB_TEST_ENV="bats-core"
if [ "${0#*$ORB_TEST_ENV}" == "$0" ]; then
    AddMissingChgLogConfig
fi
