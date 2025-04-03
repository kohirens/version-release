#!/bin/sh

set -e

find /tmp/testdata -type f -print0 | while IFS= read -r -d $'\0' file; do
  echo "cloning ${file}"
  git clone "${file}"
done
