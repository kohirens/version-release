#!/bin/sh

set -e

getopt --test > /dev/null && true
if [ $? -ne 4 ]; then
    echo 'sorry, getopts --test` failed in this environment'
    exit 1
fi

LONG_OPTS=build-arg:,entrypoint:,name:,password-stdin,progress:,rm,target:
OPTIONS=c:,d,f:,t:,u:

PARSED=$(getopt --options=${OPTIONS} --longoptions=${LONG_OPTS} --name "$0" -- "${@}") || exit 1
eval set -- "${PARSED}"

while true; do
    case "${1}" in
        -c)
            shift 2
            ;;
        -d)
            shift
            ;;
        -f)
            build_file="${2}"
            shift 2
            ;;
        -t)
            tag_name=${2}
            shift 2
            ;;
        -u)
            shift 2
            ;;
        --build-arg)
            build_args="${2}"
            shift 2
            ;;
        --entrypoint)
            build_args="${2}"
            shift 2
            ;;
        --name)
            shift 2
            ;;
        --password-stdin)
            read password_str
            shift
            ;;
        --progress)
            shift 2
            ;;
        --rm)
            shift
            ;;
        --target)
            build_target="${2}"
            shift 2
            ;;
        --) shift; break;;
#        *) echo "unknown option '${1}'"; exit 1;;
    esac
done


if [ "$#" -lt 1 ]; then
    echo "the subcommand is a required first argument"
    exit 1
fi

if [ "${1}" = "build" ] && [ "$#" -lt 2 ]; then
    echo "missing required second argument conext for the build command"
    exit 1
fi

if [ "${1}" = "build" ]; then
    build_context="${2}"

    echo "checking docker build command input"
    if [ "${tag_name}" != "${e_tag_name}" ]; then
        echo "expected tag ${e_tag_name}, got ${tag_name}"
        exit 1
    fi

    if [ "${build_file}" != "${e_build_file}" ]; then
        echo "expected file ${e_build_file}, got ${build_file}"
        exit 1
    fi

    if [ "${build_target}" != "${e_build_target}" ]; then
        echo "expected target ${e_build_target}, got ${build_target}"
        exit 1
    fi

    if [ "${build_args}" != "${e_build_args}" ]; then
        echo "expected build args ${e_build_args}, got ${build_args}"
        exit 1
    fi

    if [ "${build_context}" != "${e_build_context}" ]; then
        echo "expected context ${e_build_context}, got ${build_context}"
        exit 1
    fi
fi


if [ "${1}" = "tag" ]; then
    image="${2}"
    tag="${3}"

    echo "checking docker tag command input"
    if [ "${image}" != "${e_image}" ]; then
        echo "expected image ${e_image}, got ${image}"
        exit 1
    fi

    if [ "${tag}" != "example/com:latest" ] && [ "${tag}" != "example/com:1.0.0" ]; then
        echo "expected tags latest or 1.0.0, got ${tag}"
        exit 1
    fi
fi

if [ "${1}" = "push" ]; then
    image="${2}"
    echo "checking docker push command input"

    if [ "${image}" != "example/com:latest" ] && [ "${image}" != "example/com:1.0.0" ]; then
        echo "expected tags latest or 1.0.0, got ${image}"
        exit 1
    fi
fi

if [ "${1}" = "rmi" ]; then
    image="${2}"
    echo "checking docker rmi command input"

    if [ "${image}" != "example/com:latest" ] && [ "${image}" != "example/com:1.0.0" ] && [ "${image}" != "example/com:abc123" ]; then
        echo "expected tags latest or 1.0.0, abc123, got ${image}"
        exit 1
    fi
fi
