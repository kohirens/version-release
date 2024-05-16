InstallPseudoDocker() {
    if [ ! -f "/home/circleci/bin/docker" ]; then
        wd=$(pwd)
        cp src/tests/bin/docker /home/circleci/bin
    fi
    chmod +x /home/circleci/bin/docker
    export PATH="/home/circleci/bin:${PATH}"
}

# Runs prior to every test
setup() {
    # Install required software
    InstallPseudoDocker
}

@test '1: publish multiple docker images' {
    # setup
    export BUILD_CONTEXT="."
    export TAGS="latest 1.0.0"
    export REPOSITORY="example/com"
    export CIRCLE_SHA1="abc123"
    export DH_PASS="fakePass"
    export DH_USER="fakeUser"
    export DOCKER_FILE="path/to/Dockerfile"
    export TARGET="test-target"
    export BUILD_ARGS="--build-arg ARG1"

    # test
    source ./src/scripts/publish-docker-hub.sh

    # assert
    result="${?}"
    [ "${result}" == "0" ]
}

@test '2: publish a docker image' {
    # setup
    export BUILD_CONTEXT="."
    export IMG_TAG="1.0.0"
    export REPOSITORY="example/com"
    export CIRCLE_SHA1="abc123"
    export DH_PASS="fakePass"
    export DH_USER="fakeUser"
    export DOCKER_FILE="path/to/Dockerfile"
    export TARGET="test-target"
    export BUILD_ARGS="--build-arg ARG1"

    # test
    source ./src/scripts/publish-docker-hub.sh

    # assert
    result="${?}"
    [ "${result}" == "0" ]
}