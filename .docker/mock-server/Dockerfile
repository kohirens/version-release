ARG REPO='github.com/kohirens/version-release'
ARG GIT_CLIFF_VER="2.2.1"

FROM golang:1.23-alpine3.20 AS base

WORKDIR /tmp

ARG REPO
ARG GIT_CLIFF_VER

ENV GOPATH /root
ENV CGO_ENABLED=0
ENV SHELL=/bin/sh
ENV TZ=America/Detroit

RUN echo "https://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories  \
 && apk --no-progress --purge --no-cache upgrade \
 && apk --no-progress --purge --no-cache add --upgrade \
    curl \
    git \
    git-cliff \
    git-daemon \
    openssh \
    openssl \
    tzdata \
 && apk --no-progress --purge --no-cache upgrade \
 && rm -vrf /var/cache/apk/* \
 && cp /usr/share/zoneinfo/America/Detroit /etc/localtime

HEALTHCHECK --interval=60s --timeout=3s \
  CMD curl -f https://127.0.0.1/ || exit 1

COPY --chmod=755 .docker/mock-server/start.sh /usr/local/bin

RUN chmod 755 /usr/local/bin/start.sh

ENTRYPOINT [ "start.sh" ]

CMD []

# Make directories with current user permissions
RUN mkdir -p ~/bin ~/src

ENV PATH=${PATH}:/root/bin

COPY --chmod=755 .docker/mock-server/gen-ss-cert.sh /usr/local/bin/
COPY .docker/.gitconfig /root

RUN chmod 755 /usr/local/bin/gen-ss-cert.sh \
 && gen-ss-cert.sh --company="mock-server" \
    --sans="DNS:github.com, DNS:*.github.com, DNS:*.circleci.com" \
    --out-dir="./ssl" \
    "circleci.com"

COPY avr/testdata/repo-01.bundle /tmp/testdata/repo-01.bundle
COPY avr /root/src/${REPO}/avr/

# Go build rules require the repo to be in the GOPATH
# We also need to have docker watch for changes to any file in the avr directory
# so that it can rebuild and restart the mock-server on file changes.
WORKDIR /root/src/${REPO}/avr

RUN cd mock-server \
 && mkdir -p request tmp \
 && go mod tidy \
 && go generate \
 && go build \
 && go install
