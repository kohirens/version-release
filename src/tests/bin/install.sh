#!/bin/sh

set -e

InstallPseudoDocker() {
    if [ ! -f "/home/circleci/bin/docker" ]; then
        wd=$(pwd)
        cp src/tests/data/docker /home/circleci/bin
    fi
}

InstallPseudoCmd() {
    c_m_d="${1}"
    if [ ! -f "/home/circleci/bin/${c_m_d}" ]; then
        wd=$(pwd)
        cp src/tests/bin/"${c_m_d}".sh /home/circleci/bin/"${c_m_d}"
    fi
    chmod +x /home/circleci/bin/${c_m_d}
    export PATH="/home/circleci/bin:${PATH}"
}

export PATH="/home/circleci/bin:${PATH}"
