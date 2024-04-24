ARG REPO='github.com/kohirens/version-release-orb'
ARG GIT_CLIFF_VER="2.2.1"

FROM golang:1.21-alpine3.19 AS base

ARG REPO
ARG GIT_CLIFF_VER

ENV GOPATH /root
ENV CGO_ENABLED=0
ENV SHELL=/bin/sh
ENV TZ=America/Detroit

RUN env && apk --no-progress --purge --no-cache upgrade \
 && apk --no-progress --purge --no-cache add --upgrade \
    curl \
    git \
    git-daemon \
    openssh \
    openssl \
    tzdata \
 && apk --no-progress --purge --no-cache upgrade \
 && rm -vrf /var/cache/apk/* \
 && cp /usr/share/zoneinfo/America/Detroit /etc/localtime

# Make directories with current user permissions
RUN mkdir -p ~/bin ~/src

ENV PATH=${PATH}:/root/bin

# Install git-cliff
RUN wget https://github.com/orhun/git-cliff/releases/download/v${GIT_CLIFF_VER}/git-cliff-${GIT_CLIFF_VER}-aarch64-unknown-linux-musl.tar.gz \
 && tar -xz -f git-cliff-${GIT_CLIFF_VER}-aarch64-unknown-linux-musl.tar.gz \
 && mv git-cliff-${GIT_CLIFF_VER}/git-cliff /usr/local/bin \
 && git-cliff --version

COPY --from=kohirens/git-tool-belt:2.1.2 /usr/local/bin/git-tool-belt /root/bin

COPY --chmod=+x .docker/mock-server/gen-ss-cert.sh /usr/bin/
COPY --chown=root:root .docker/.gitconfig /root

RUN /usr/bin/gen-ss-cert.sh --company="mock-server" \
    --sans="DNS:github.com, DNS:*.github.com, DNS:*.circleci.com" \
    --out-dir="./ssl" \
    "circleci.com"

HEALTHCHECK --interval=5s --timeout=3s \
  CMD curl -f https://127.0.0.1/ || exit 1

WORKDIR /tmp

COPY vro/testdata/repo-01.bundle ./

WORKDIR "/root/src/${REPO}/vro"

COPY ./vro/ ./

# Build the mock server app
RUN cd mock-server \
 && mkdir -p request tmp \
 && go mod tidy \
 && go generate \
 && go build \
 && go install

COPY --chmod=+x .docker/mock-server/start.sh /usr/local/bin

ENTRYPOINT [ "start.sh" ]

CMD []