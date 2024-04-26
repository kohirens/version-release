ARG REPO='github.com/kohirens/version-release-orb'
ARG GIT_CLIFF_VER="2.2.1"

FROM golang:1.21-alpine3.19 AS base

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

# Make directories with current user permissions
RUN mkdir -p ~/bin ~/src

ENV PATH=${PATH}:/root/bin

COPY --from=kohirens/git-tool-belt:2.1.2 /usr/local/bin/git-tool-belt /root/bin

COPY --chown=root:root --chmod=+x .docker/mock-server/gen-ss-cert.sh /usr/bin/
COPY --chown=root:root .docker/.gitconfig /root

RUN chmod +x /usr/bin/gen-ss-cert.sh \
 && gen-ss-cert.sh --company="mock-server" \
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

RUN chmod +x /usr/local/bin/start.sh

ENTRYPOINT [ "start.sh" ]

CMD []