#!/bin/sh

# We'll use this script to manage starting and stopping this container gracefully.
# It only takes up about 00.01 CPU % allotted to the container, you can verify
# by running `docker stats` after you start a container that uses this as
# as the CMD.

set -e

shutd () {
    printf "Shutting down the container gracefully..."
    # You can run clean commands here!
    last_signal="15"
}

trap 'shutd' TERM

echo "Starting up..."

OUT_DIR="mock-server/ssl"

mkdir -p "${OUT_DIR}"

# Run non-blocking commands here
gen-ss-cert.sh --company="mock-server" \
    --sans="DNS:github.com, DNS:*.github.com, DNS:*.circleci.com, IP:127.0.0.1, IP:0:0:0:0:0:0:0:1" \
    --out-dir="${OUT_DIR}" \
    "circleci.com"

mkdir -p mock-server/tmp

# Run the mock server (git-web) in the background
mock-server &
child_proc=${!}

echo "Ready!"

# This keeps the container running until it receives a signal to be stopped.
# Also very low CPU usage.
if [ "${KEEP_RUNNING}" = "1" ]; then
    while [ "${last_signal}" != "15" ]; do
        sleep 1
    done
else
    go test ./...
fi

# kill the mock server
kill -9 ${child_proc}
wait ${child_proc}
echo "done"
