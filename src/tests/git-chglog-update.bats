# Runs prior to every test
setup() {
    wd=$(pwd)
    curl -LO https://github.com/git-chglog/git-chglog/releases/download/v0.15.1/git-chglog_0.15.1_linux_amd64.tar.gz
    mv git-chglog_0.15.1_linux_amd64.tar.gz /tmp/git-chglog_0.15.1_linux_amd64.tar.gz
    cd /tmp
    tar -xzf git-chglog_0.15.1_linux_amd64.tar.gz
    chmod +x ./git-chglog
    mkdir -p /home/circleci/bin
    mv ./git-chglog /home/circleci/bin
    cd "${wd}"
    # Load our script file.
    source ./src/scripts/git-chglog-update.sh
}

@test '1: Update CHANGELOG' {
    # Mock environment variables or functions by exporting them (after the script has been sourced)
    export PARAM_CONFIGFILE="src/tests/data/.chglog/config.yml"
    export PARAM_OUTPUTFILE="CHANGLOG-3000.md"
    # Run the function
    GitChglogUpdate
    # Check that the output file was produces.
    [ -f "${PARAM_OUTPUTFILE}" ]
    cat $PARAM_OUTPUTFILE
    #rm $PARAM_OUTPUTFILE
}