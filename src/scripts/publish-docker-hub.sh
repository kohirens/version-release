if [ -z "${BUILD_CONTEXT}" ]; then
    echo "cannot build the image, build_context parameter is empty"
    exit 1
fi

# Get the value of the semantic version tag in 1 of 3 way.
# If more than 1 is set, the last wins.
semver=""
if [ -n "${SEMVER_TAG_CMD}" ]; then
    semver="$("${SEMVER_TAG_CMD}")"
    echo "semantic version ${semver} was set by command"
fi

if [ -n "${SEMVER_TAG_ENV_VAR}" ]; then
    semver="${!SEMVER_TAG_ENV_VAR}"
    echo "semantic version ${semver} was extracted from environment variable ${SEMVER_TAG_ENV_VAR}"
fi

if [ -n "${SEMVER_TAG_FILE}" ]; then
    semver="$(cat "${SEMVER_TAG_FILE}")"
    echo "semantic version ${semver} was pulled from the file ${SEMVER_TAG_FILE}"
fi

if [ -n "${SEMVER_TAG}" ]; then
    semver="${SEMVER_TAG}"
    echo "semantic version ${semver} was explicitly set"
fi

if [ -z "${semver}" ]; then
    echo "cannot build the image because semver_tag, semver_tag_cmd, semver_tag_env_var, and semver_tag_file parameters are all empty; at least 1 needs to be set"
    exit 1
fi

if [ -z "${REPOSITORY}" ]; then
    echo "cannot build the image, repository parameter is empty"
    exit 1
fi

if [ -z "${DH_API_TOKEN}" ]; then
    echo "the environment variable DH_API_TOKEN containing the Docker Hub API token/password is not set"
    exit 1
fi

if [ -n "${ENV_FILE}" ]; then
    # shellcheck disable=SC1090
    source "${ENV_FILE}"
fi

export DH_IMAGE="${REPOSITORY}:${CIRCLE_SHA1}"

echo "${DH_API_TOKEN}" | docker login -u "${DH_USER}" --password-stdin

build_cmd="docker build --rm -t ${DH_IMAGE} -f ${DOCKER_FILE}"

if [ -n "${TARGET}" ]; then
    build_cmd="${build_cmd} --target ${TARGET}"
fi

if [ -n "${BUILD_ARGS}" ]; then
    build_cmd="${build_cmd} ${BUILD_ARGS}"
fi

# Add the version LABEL
if [ -n "${semver}" ]; then
    echo "adding build argument semver to build command"
    build_cmd="${build_cmd} --build-arg SEMVER=${semver}"
fi

build_cmd="${build_cmd} ${BUILD_CONTEXT}"

echo "executing build command:"
printf "\t%s\n" "${build_cmd}"

$build_cmd

# The user should not prefix the semver tag since it is done here automatically.
# If they do, it will just perform an overwrite, which could throw an
# error if they do not have permissions set appropriately.
TAGS="${semver} ${TAGS}"

if [ -n "${TAGS}" ]; then
    printf "additional tags found: %s\n" "${TAGS}"
fi

for tag in ${TAGS}; do
    stamp="${REPOSITORY}:${tag}"

    docker tag "${DH_IMAGE}" "${stamp}"

    printf "pushing %s\n" "${stamp}"
    docker push "${stamp}"

    printf "cleaning up %s\n" "${stamp}"
    docker rmi "${stamp}"
done

printf "cleaning up %s\n" "${DH_IMAGE}"
docker rmi "${DH_IMAGE}"
