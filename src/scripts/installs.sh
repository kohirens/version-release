InstallGitChglog() {
    if [ ! -f "/home/circleci/bin/git-chglog" ]; then
        wd=$(pwd)
        curl -LO https://github.com/git-chglog/git-chglog/releases/download/v0.15.1/git-chglog_0.15.1_linux_amd64.tar.gz
        mv git-chglog_0.15.1_linux_amd64.tar.gz /tmp/git-chglog_0.15.1_linux_amd64.tar.gz
        cd /tmp || exit 1
        tar -xzf git-chglog_0.15.1_linux_amd64.tar.gz
        chmod +x ./git-chglog
        mkdir -p /home/circleci/bin
        mv ./git-chglog /home/circleci/bin
        cd "${wd}" || exit 1
    fi
}

InstallGitToolBelt() {
    if [ ! -f "/home/circleci/bin/git-tool-belt" ]; then
        wd=$(pwd)
        curl -LO https://github.com/kohirens/git-tool-belt/releases/download/2.1.0/git-tool-belt-linux-amd64.tar.gz
        mv git-tool-belt-linux-amd64.tar.gz /tmp/git-tool-belt-linux-amd64.tar.gz
        cd /tmp || exit 1
        tar -xzf git-tool-belt-linux-amd64.tar.gz
        chmod +x ./git-tool-belt-linux-amd64
        mkdir -p /home/circleci/bin
        mv ./git-tool-belt-linux-amd64 /home/circleci/bin/git-tool-belt
        cd "${wd}" || exit 1
    fi
}
