ARG USER_NAME='circleci'
ARG USER_UID='1000'
ARG USER_GID='1000'
ARG USER_GROUP='app_users'
ARG REPO='github.com/kohirens/version-release-orb'
ARG GIT_CLIFF_VER="2.2.1"
FROM kohirens/go:1.21.0 AS build

ARG USER_NAME
ARG USER_UID
ARG USER_GID
ARG USER_GROUP
ARG REPO
ARG GIT_CLIFF_VER

ENV GOPATH /home/${USER_NAME}
ENV CGO_ENABLED=0
ENV WORK_DIR=/home/${USER_NAME}/src/${REPO}
ENV SHELL=/bin/sh

# Update OS
RUN apk --no-progress --purge --no-cache upgrade \
 && apk add \
    git \
    git-cliff \
 && rm -vrf /var/cache/apk/* \
 && rm -rf /tmp/*

RUN wget https://github.com/orhun/git-cliff/releases/download/v${GIT_CLIFF_VER}/git-cliff-${GIT_CLIFF_VER}-aarch64-unknown-linux-musl.tar.gz \
 && tar -xz -f git-cliff-${GIT_CLIFF_VER}-aarch64-unknown-linux-musl.tar.gz \
 && mv git-cliff-${GIT_CLIFF_VER}/git-cliff /usr/local/bin \
 && git-cliff --version

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

WORKDIR "/home/${USER_NAME}/src/${REPO}/vro"

# Add source files
COPY --chown=${USER_NAME}:${USER_GROUP} ./vro/ ./

# Build
RUN go mod tidy \
 && go generate \
 && go install

COPY --chmod=+x .docker/vr/start.sh /usr/local/bin
COPY --chown=${USER_NAME}:${USER_NAME} .docker/.gitconfig /home/${USER_NAME}

ENTRYPOINT [ "start.sh" ]

HEALTHCHECK --interval=5s --timeout=3s --start-period=3s --retries=4 \
    CMD vro -help || exit 1

CMD [ ]

FROM kohirens/git-tool-belt:2.1.2 AS release

ARG USER_NAME
ARG USER_GROUP
ARG REPO

HEALTHCHECK --interval=5s --timeout=3s --start-period=3s --retries=4 \
    CMD vro -help || exit 1

USER root

# Install git-cliff
RUN wget https://github.com/orhun/git-cliff/releases/download/v${GIT_CLIFF_VER}/git-cliff-${GIT_CLIFF_VER}-aarch64-unknown-linux-musl.tar.gz \
 && tar -xz -f git-cliff-${GIT_CLIFF_VER}-aarch64-unknown-linux-musl.tar.gz \
 && mv git-cliff-${GIT_CLIFF_VER}/git-cliff /usr/local/bin \
 && git-cliff --version

RUN rm /usr/local/bin/git-chglog

USER ${USER_NAME}

# Add build artifacts
COPY --from=build --chown=${USER_NAME}:${USER_GROUP} "/home/${USER_NAME}/bin/vro" /bin/vro
COPY --from=build --chown=${USER_NAME}:${USER_GROUP} "/home/${USER_NAME}/bin/git-chglog" /bin

ENV PATH="${PATH}:/home/${USER_NAME}/bin"

ENTRYPOINT [ "vro" ]
