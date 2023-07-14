#!/bin/sh

set -e

if [ -z "${BUILD_CONTEXT}" ]; then
    echo "cannot build the image, build_context parameter is an empty"
    exit 1
fi

if [ -z "${IMG_TAG}" ]; then
    IMG_TAG=$(git describe --tags --always)
    echo "image tag parameter is empty, defaulting to ${IMG_TAG}"
fi

export DH_IMAGE="${REPOSITORY}:${IMG_TAG}"

echo "${DH_PASS}" | docker login -u "${DH_USER}" --password-stdin

build_cmd="docker build --rm -t ${DH_IMAGE} -f ${DOCKER_FILE}"

if [ -n "${TARGET}" ]; then
    build_cmd="${build_cmd} --target ${TARGET}"
fi

if [ -n "${BUILD_ARGS}" ]; then
    build_cmd="${build_cmd} ${BUILD_ARGS}"
fi

build_cmd="${build_cmd} ${BUILD_CONTEXT}"

printf "\nBuilding %s\n" "${DH_IMAGE}"
$build_cmd

printf "\nPushing %s\n" "${DH_IMAGE}"
docker push "${DH_IMAGE}"

printf "\nCleaning up %s\n" "${DH_IMAGE}"
docker rmi "${DH_IMAGE}"
