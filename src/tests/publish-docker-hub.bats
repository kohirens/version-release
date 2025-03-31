# Runs prior to every test
setup() {
    source ./src/tests/bin/install.sh
    InstallPseudoCmd "docker"
    # Install required software
}

#share input
export BUILD_CONTEXT="."
export REPOSITORY="example/com"
export CIRCLE_SHA1="abc123"
export DH_API_TOKEN="fakePass"
export DH_USER="fakeUser"
export DOCKER_FILE="path/to/Dockerfile"
export TARGET="test-target"
export BUILD_ARGS="--build-arg ARG1"

# Set expected arguments for the mod docker command.
export e_tag_name="example/com:abc123"
export e_build_file="path/to/Dockerfile"
export e_build_target="test-target"
export e_build_args="ARG1"
export e_build_context="."
export e_image="example/com:abc123"

@test '1: publish multiple docker images' {
    # setup
    export TAGS="latest 1.0.0"

    # test
    source ./src/scripts/publish-docker-hub.sh

    # assert
    result="${?}"
    [ "${result}" == "0" ]
}

@test '2: publish a docker image' {
    # setup
    export IMG_TAG="1.0.0"

    # test
    source ./src/scripts/publish-docker-hub.sh

    # assert
    result="${?}"
    [ "${result}" == "0" ]
}