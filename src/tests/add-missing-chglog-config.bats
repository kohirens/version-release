# Runs prior to every test
setup() {
    wd=$(pwd)
    curl -LO https://github.com/kohirens/git-tool-belt/releases/download/1.2.2/git-tool-belt-linux-amd64.tar.gz
    mv git-tool-belt-linux-amd64.tar.gz /tmp/git-tool-belt-linux-amd64.tar.gz
    cd /tmp
    tar -xzf git-tool-belt-linux-amd64.tar.gz
    chmod +x ./git-tool-belt-linux-amd64
    mkdir -p /home/circleci/bin
    mv ./git-tool-belt-linux-amd64 /home/circleci/bin/git-tool-belt
    cd "${wd}"
    # Load our script file.
    source ./src/scripts/add-missing-chglog-config.sh
}

@test '1: Check for git-chglog config' {
    # Mock environment variables or functions by exporting them (after the script has been sourced)
    export PARAM_CONFIG_FILE="./test/.chglog/config.yml"
    export CIRCLE_REPOSITORY_URL="1234"
    # Run the function
    AddMissingChgLogConfig
    # Capture the return value of the function
    [ -f "${PARAM_CONFIG_FILE}" ];
}