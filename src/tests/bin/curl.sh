#!/bin/sh

set -e

getopt --test > /dev/null && true
if [ $? -ne 4 ]; then
    echo 'sorry, getopts --test` failed in this environment'
    exit 1
fi

LONG_OPTS=url:,show-error,silent,output:
OPTIONS=H:,k,L,w:

PARSED=$(getopt --options=${OPTIONS} --longoptions=${LONG_OPTS} --name "$0" -- "${@}") || exit 1
eval set -- "${PARSED}"

url=""

while true; do
    case "${1}" in
        --url) url="${2}"
            shift 2
            ;;
        --) shift; break;;
        *) shift;;
    esac
done

case "${url}" in
    https://api.github.com/repos/kohirens/version-release/releases/tags/4.2.4)
        echo "200";
        ;;
    https://api.github.com/repos/kohirens/version-release/releases/tags/f4.2.4-fake)
        echo "404";
        ;;
        *) echo "skipping \"${url}\""; exit 1;;
esac

exit 0
