if [ -z "${BUILD_CONTEXT}" ]; then
    echo "cannot build the image, build_context parameter is empty"
    exit 1
fi

if [ -z "${IMG_TAG}" ] && [ -z "${TAGS}" ]; then
    echo "cannot build the image, image_tag and tags parameters are empty"
    exit 1
fi

# Get the value of the semantic version tag in 1 of 3 way.
# If more than 1 is set, the last wins.
semver=""
if [ -n "${PARAM_TAG_CMD}" ]; then
    semver="$("${PARAM_TAG_CMD}")"
    echo "semantic version ${semver} was set by command"
fi

if [ -n "${PARAM_TAG_ENV_VAR}" ]; then
    semver="${!PARAM_TAG_ENV_VAR}"
    echo "semantic version ${semver} was extracted from environment variable ${PARAM_TAG_ENV_VAR}"
fi

if [ -n "${PARAM_TAG_FILE}" ]; then
    semver="$(cat "${PARAM_TAG_FILE}")"
    echo "semantic version ${semver} was pull from file ${PARAM_TAG_FILE}"
fi

# For backward compatibility.
if [ -n "${semver}" ]; then
    TAGS="${semver}"
elif [ -n "${IMG_TAG}" ] && [ -z "${TAGS}" ]; then
    TAGS="${IMG_TAG}"
elif [ -n "${IMG_TAG}" ] && [ -n "${TAGS}" ]; then
    TAGS="${TAGS} ${IMG_TAG}"
fi

if [ -z "${REPOSITORY}" ]; then
    echo "cannot build the image, repository parameter is empty"
    exit 1
fi

if [ -n "${ENV_FILE}" ]; then
    ./"${ENV_FILE}"
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
