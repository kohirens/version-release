ARG USER_NAME='app'
ARG USER_UID='1000'
ARG USER_GID='1000'
ARG USER_GROUP='app_users'
ARG REPO='github.com/kohirens/version-release'
ARG GIT_CLIFF_VER="2.8.0"
ARG SEMVER

FROM kohirens/go:1.21.0 AS build

ARG USER_NAME
ARG USER_UID
ARG USER_GID
ARG USER_GROUP
ARG REPO
ARG GIT_CLIFF_VER

ENV GOPATH /home/${USER_NAME}
ENV CGO_ENABLED=0
ENV SHELL=/bin/sh
ENV TZ=America/Detroit

# Update OS
RUN echo "https://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories \
 && apk --no-progress --purge --no-cache upgrade \
 && apk add \
    curl \
    git \
 && rm -vrf /var/cache/apk/* \
 && rm -rf /tmp/*

RUN wget https://github.com/orhun/git-cliff/releases/download/v${GIT_CLIFF_VER}/git-cliff-${GIT_CLIFF_VER}-x86_64-unknown-linux-musl.tar.gz \
 && tar -xz -f git-cliff-${GIT_CLIFF_VER}-x86_64-unknown-linux-musl.tar.gz \
 && git-cliff-${GIT_CLIFF_VER}/git-cliff --version \
 && mv git-cliff-${GIT_CLIFF_VER}/git-cliff /usr/local/bin \
 && chmod 755 /usr/local/bin/git-cliff \
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

# Build
# Add source files
COPY --chown=${USER_NAME}:${USER_GROUP} ./avr/ ./avr

RUN cd avr \
 && go mod tidy \
 && go generate \
 && go install

COPY --chmod=755 .docker/vr/start.sh /usr/local/bin
COPY --chown=${USER_NAME}:${USER_NAME} .docker/.gitconfig /home/${USER_NAME}

ENTRYPOINT [ "start.sh" ]

HEALTHCHECK --interval=5s --timeout=3s --start-period=3s --retries=4 \
    CMD avr -help || exit 1

CMD [ ]

FROM kohirens/git-tool-belt:2.1.2 AS release

ARG USER_NAME
ARG USER_GROUP
ARG REPO
ARG GIT_CLIFF_VER

HEALTHCHECK --interval=5s --timeout=3s --start-period=3s --retries=4 \
    CMD avr -help || exit 1

USER root

RUN echo "https://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories \
 && apk --no-progress --purge --no-cache upgrade \
 && apk add \
    curl \
 && rm -vrf /var/cache/apk/* \
 && rm -rf /tmp/*

RUN wget https://github.com/orhun/git-cliff/releases/download/v${GIT_CLIFF_VER}/git-cliff-${GIT_CLIFF_VER}-x86_64-unknown-linux-musl.tar.gz \
 && tar -xz -f git-cliff-${GIT_CLIFF_VER}-x86_64-unknown-linux-musl.tar.gz \
 && git-cliff-${GIT_CLIFF_VER}/git-cliff --version \
 && mv git-cliff-${GIT_CLIFF_VER}/git-cliff /usr/local/bin \
 && chmod 755  /usr/local/bin/git-cliff \
 && git-cliff --version

RUN rm /usr/local/bin/git-chglog

# Add build artifacts
COPY --link --from=build "/home/${USER_NAME}/bin/avr" "/usr/local/bin/avr"
COPY --link --chmod=755 ./src/scripts/ /usr/local/bin/

# DEBUG
RUN ls -la /usr/local/bin \
 && chmod 755 /usr/local/bin/*.sh \
 && ls -la /usr/local/bin

USER ${USER_NAME}

# Flatten the release image to a few layers
FROM scratch AS flattened-release

ARG SEMVER

COPY --from=release / /

ENTRYPOINT [ "avr" ]

CMD [ "--help" ]

# docker build --rm -t kohirens/version-release:dev -f .docker/vr/Dockerfile .

LABEL "com.kohirens.vendor"="Kohirens LLC"
LABEL version="${SEMVER}"
LABEL description="Auto Version Release. \
Automatically update your CHANGELOG and a release on GitHub using."