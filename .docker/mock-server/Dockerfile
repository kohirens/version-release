ARG REPO='github.com/kohirens/version-release-orb'

FROM golang:1.20-alpine3.18 AS base

ARG REPO

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

#RUN go install github.com/git-chglog/git-chglog/cmd/git-chglog@v0.15.4 \
# && go clean -cache -modcache -testcache

RUN cd /tmp \
 && git clone https://github.com/b01/git-chglog.git \
 && git config --global --add safe.directory /tmp/git-chglog/git-chglog \
 && cd git-chglog \
 && git checkout -t origin/issue-76-enhancement-when-no-tags-present \
 && export VERSION=$(git describe --tags --always) \
 && go install -ldflags="-X 'main.version=${VERSION}'" ./cmd/git-chglog \
 && go clean -cache -modcache -testcache

COPY --from=kohirens/git-tool-belt:2.1.2 /usr/local/bin/git-tool-belt /root/bin

COPY --chmod=+x .docker/mock-server/gen-ss-cert.sh /usr/bin/
COPY --chmod=+x .docker/mock-server/start.sh /usr/bin/

RUN chmod +x /usr/bin/gen-ss-cert.sh
RUN chmod +x /usr/bin/start.sh
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

# Build
RUN cd mock-server \
 && mkdir -p request tmp \
 && go mod tidy \
 && go generate \
 && go build \
 && go install

ENV PATH=${PATH}:/root/bin

ENTRYPOINT [ "/usr/bin/start.sh" ]

CMD []