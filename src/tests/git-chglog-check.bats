# Runs prior to every test
setup() {
    # Load our script file.
    source ./src/scripts/git-chglog-check.sh
}

@test '1: Init check for git-chglog' {
    # Mock environment variables or functions by exporting them (after the script has been sourced)
    export PARAM_CONFIGFILE="./test/.chglog/config.yml"
    # Capture the output of our "Greet" function
    result=$(Greet)
    [ "$result" == "0" ]
}