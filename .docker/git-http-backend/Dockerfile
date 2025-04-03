ARG REPO='github.com/kohirens/version-release'

FROM httpd:2.4.63-alpine AS base

WORKDIR /tmp

ARG REPO

ENV SHELL=/bin/sh
ENV TZ=America/Detroit

RUN echo "https://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories  \
 && apk --no-progress --purge --no-cache upgrade \
 && apk --no-progress --purge --no-cache add --upgrade \
    apache2 \
    curl \
    git \
    git-daemon \
    openssh \
    openssl \
    tzdata \
 && apk --no-progress --purge --no-cache upgrade \
 && rm -vrf /var/cache/apk/* \
 && cp /usr/share/zoneinfo/America/Detroit /etc/localtime

HEALTHCHECK --interval=120s --timeout=3s CMD curl -f https://127.0.0.1/ || exit 1

COPY .docker/.gitconfig /root

#RUN curl -L -o gen-ss-cert.sh https://raw.githubusercontent.com/b01/script-lib/refs/heads/main/generate-self-signed-cert.sh \
# && chmod +x ./gen-ss-cert.sh \
# && mv gen-ss-cert.sh /usr/local/bin

COPY .docker/git-http-backend/my-httpd.conf /usr/local/apache2/conf/httpd.conf
COPY .docker/git-http-backend/my-httpd-ssl.conf /usr/local/apache2/conf/extra/httpd-ssl.conf
COPY avr/testdata/*.bundle /tmp/testdata/
COPY avr/pkg/github/testdata/*.bundle /tmp/testdata/
COPY --chmod=755 .docker/git-http-backend/expand-bundles.sh /tmp

# Copy the SSL cert into the container where Apache expects to find them.
COPY --chmod=755 --chown=www-data:www-data .docker/ssl/certs/server.pem /usr/local/apache2/conf/server.crt
COPY --chmod=755 --chown=www-data:www-data .docker/ssl/private/server.key /usr/local/apache2/conf/server.key
COPY --chmod=755 --chown=www-data:www-data .docker/git-http-backend/passwords /usr/local/apache2/conf/passwords

RUN mkdir -p /usr/local/apache2/htdocs/kohirens

WORKDIR /usr/local/apache2/htdocs/kohirens

RUN /tmp/expand-bundles.sh

RUN chown -R www-data:www-data /usr/local/apache2/htdocs/kohirens
RUN cat /usr/local/apache2/conf/server.crt >> /etc/ssl/certs/ca-certificates.crt \
 && echo "" >> /etc/ssl/certs/ca-certificates.crt
