ARG USER_NAME='circleci'
ARG USER_UID='1000'
ARG USER_GID='1000'
ARG USER_GROUP='app_users'
ARG REPO='github.com/kohirens/version-release-orb'

FROM kohirens/go:1.21.0 AS build

ARG USER_NAME
ARG USER_UID
ARG USER_GID
ARG USER_GROUP
ARG REPO

ENV GOPATH /home/${USER_NAME}
ENV CGO_ENABLED=0
ENV WORK_DIR=/home/${USER_NAME}/src/${REPO}
ENV SHELL=/bin/sh

# Update OS
RUN apk --no-progress --purge --no-cache upgrade \
 && apk add \
    git \
 && rm -vrf /var/cache/apk/* \
 && rm -rf /tmp/*

# Add a non-root group and user.
RUN addgroup --system --gid ${USER_GID} ${USER_GROUP} \
 && adduser --system \
    --disabled-password \
    --ingroup ${USER_GROUP} \
    --uid ${USER_UID} \
    ${USER_NAME}

USER ${USER_NAME}

# Make directories with current user permissions
RUN mkdir -p ~/bin ~/src

ENV PATH="${PATH}:/home/${USER_NAME}/bin"

COPY --from=kohirens/git-tool-belt:2.1.2 --chown=${USER_NAME}:${USER_GROUP} /usr/local/bin/git-tool-belt /home/${USER_NAME}/bin

# Don't install latest until this patch is merged.
#RUN go install -ldflags="-X 'main.version=${VERSION}'" github.com/git-chglog/git-chglog/cmd/git-chglog@latest \
# && go clean -cache -modcache -testcache

# Install patched forked version with the requied fix for now.
RUN cd /tmp \
 && git clone https://github.com/b01/git-chglog.git \
 && git config --global --add safe.directory /tmp/git-chglog/git-chglog \
 && cd git-chglog \
 && git checkout -t origin/issue-76-enhancement-when-no-tags-present \
 && export VERSION=$(git describe --tags --always) \
 && go install -ldflags="-X 'main.version=${VERSION}'" ./cmd/git-chglog \
 && go clean -cache -modcache -testcache \
 && VERSION=$(git describe --tags --always) \
 && actual=$(git-chglog --version) \
 && expected="git-chglog version ${VERSION}"  \
 && echo "VERSION = \"${VERSION}\"" \
 && echo "actual = \"${actual}\"" \
 && echo "expected = \"${expected}\"" \
 && [ "${actual}" = "${expected}" ]

WORKDIR "/home/${USER_NAME}/src/${REPO}"

# Add source files
COPY --chown=${USER_NAME}:${USER_GROUP} ./vro/ ./

# Build
RUN go mod tidy \
 && go generate \
 && go install

COPY .docker/vr/start.sh /usr/local/bin/start.sh

ENTRYPOINT [ "start.sh" ]

HEALTHCHECK --interval=5s --timeout=3s --start-period=3s --retries=4 \
    CMD vro -help || exit 1

USER root

CMD [ ]

FROM kohirens/git-tool-belt:2.1.2 AS release

ARG USER_NAME
ARG USER_GROUP
ARG REPO

HEALTHCHECK --interval=5s --timeout=3s --start-period=3s --retries=4 \
    CMD vro -help || exit 1

USER root

# Update OS
RUN apk --no-progress --purge --no-cache upgrade \
 && apk add \
    git-cliff \
 && rm -vrf /var/cache/apk/* \
 && rm -rf /tmp/*

RUN rm /usr/local/bin/git-chglog

USER ${USER_NAME}

# Add build artifacts
COPY --from=build --chown=${USER_NAME}:${USER_GROUP} "/home/${USER_NAME}/bin/vro" /bin/vro
COPY --from=build --chown=${USER_NAME}:${USER_GROUP} "/home/${USER_NAME}/bin/git-chglog" /bin

ENV PATH="${PATH}:/home/${USER_NAME}/bin"

ENTRYPOINT [ "vro" ]
