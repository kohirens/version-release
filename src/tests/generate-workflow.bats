# Runs prior to every test
setup() {
    # Install required software
    source ./src/scripts/installs.sh
    InstallGitChglog
    InstallGitToolBelt

    # Load our script file.
    source ./src/scripts/git-chglog-update.sh
    source ./src/scripts/generate-workflow.sh
}

UnbundleRepo() {
    git clone -b main src/tests/data/${1}.bundle "${2}"
}

@test '1: Test generated a workflow for publish changelog' {
    local fixture
    local ret_code
    local subShell
    local params

    rm -rf "tmp/${fixture}"
    fixture="repo-02"
    # Extract repo which needs to update the changelog
    UnbundleRepo "${fixture}" "tmp/${fixture}"
    cd "tmp/${fixture}"
    export PARAM_CONFIG_FILE=".chglog/config.yml"
    export PARAM_CHANGELOG_FILE="CHANGELOG.md"
    export PARAM_REPOSITORY_PATH="."
    GitChglogUpdate

    PARAM_CHANGELOG_PATH="CHANGELOG.md"
    PARAM_GENERATED_WORKFLOW_PATH=".circleci/generated-config.yml"
    PARAM_CONTINUE_PARAMETERS_PATH=".circleci/generated-parameters.json"
    # Run the function in a way that allows us to capture the exit code.
    subShell=$(GenerateAWorkflow)
    ret_code=$?
    # Test output
    params=$(<"${PARAM_CONTINUE_PARAMETERS_PATH}")
    echo
    echo
    echo "ret_code=${ret_code}"
    echo "params=${params}"
    [ "${ret_code}" = "0" ] && [ "${params}" == '{ "continued_automation": "publish-changelog" }' ]
}

@test '2: Test generated a workflow for tag and release' {
    local fixture
    local ret_code
    local subShell
    local params

    rm -rf "tmp/${fixture}"
    fixture="repo-03"
    # Extract repo which needs to update the changelog
    UnbundleRepo "${fixture}" "tmp/${fixture}"
    cd "tmp/${fixture}"

    PARAM_CHANGELOG_PATH="CHANGELOG.md"
    PARAM_GENERATED_WORKFLOW_PATH=".circleci/generated-config.yml"
    PARAM_CONTINUE_PARAMETERS_PATH=".circleci/generated-parameters.json"
    # Run the function in a way that allows us to capture the exit code.
    subShell=$(GenerateAWorkflow)
    echo "${subShell}"
    ret_code=$?
    # Test output
    params=$(<"${PARAM_CONTINUE_PARAMETERS_PATH}")
    echo
    echo
    echo "ret_code=${ret_code}"
    echo "params=${params}"
    [ "${ret_code}" = "0" ] && [ "${params}" == '{ "continued_automation": "publish-new-tag" }' ]
}