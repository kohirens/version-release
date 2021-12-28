# Runs prior to every test
setup() {
    # Load our script file.
    source ./src/scripts/git-chglog-update.sh
}

@test '1: Check for git-chglog config' {
    # Mock environment variables or functions by exporting them (after the script has been sourced)
    export PARAM_CONFIGFILE="./test/.chglog/config.yml"
    export PARAM_OUTPUTFILE="./CHANGLOG-3000.md"
    # Run the function
    GitChglogUpdate
    # Check that the output file was produces.
    [ -f "${PARAM_OUTPUTFILE}" ]
    cat $PARAM_OUTPUTFILE
    #rm $PARAM_OUTPUTFILE
}