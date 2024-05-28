#!/bin/sh

set -e

if [ -z "${BUILD_CONTEXT}" ]; then
    echo "cannot build the image, build_context parameter is empty"
    exit 1
fi

if [ -z "${IMG_TAG}" ] && [ -z "${TAGS}" ]; then
    echo "cannot build the image, image_tag and tags parameters are empty"
    exit 1
fi

# For backward compatibility.
if [ -n "${IMG_TAG}" ] && [ -z "${TAGS}" ]; then
    TAGS="${IMG_TAG}"
elif [ -n "${IMG_TAG}" ] && [ -n "${TAGS}" ]; then
    TAGS="${TAGS} ${IMG_TAG}"
fi

if [ -z "${REPOSITORY}" ]; then
    echo "cannot build the image, repository parameter is empty"
    exit 1
fi

if [ -n "${ENV_FILE}" ]; then
    export "$(cat "${ENV_FILE}")"
fi

export DH_IMAGE="${REPOSITORY}:${CIRCLE_SHA1}"

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

for tag in ${TAGS}; do
    stamp="${REPOSITORY}:${tag}"
    docker tag "${DH_IMAGE}" "${stamp}"

    printf "\nPushing %s\n" "${stamp}"
    docker push "${stamp}"

    printf "\nCleaning up %s\n" "${stamp}"
    docker rmi "${stamp}"
done

printf "\nCleaning up %s\n" "${DH_IMAGE}"
docker rmi "${DH_IMAGE}"
