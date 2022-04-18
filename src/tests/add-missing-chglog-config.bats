# Runs prior to every test
setup() {
    # Load our script file.
    source ./src/scripts/add-missing-chglog-config.sh
}

@test '1: Check for git-chglog config' {
    # Mock environment variables or functions by exporting them (after the script has been sourced)
    export PARAM_CONFIGFILE="./test/.chglog/config.yml"
    export CIRCLE_REPOSITORY_URL="1234"
    # Run the function
    AddMissingChgLogConfig
    # Capture the return value of the function
    [ -f "${PARAM_CONFIGFILE}" ];
}