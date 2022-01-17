ARG USER_NAME='app'
ARG USER_UID='1000'
ARG USER_GID='1000'
ARG USER_GROUP='app_users'
ARG REPO='github.com/kohirens/version-release-orb'

FROM kohirens/git-tool-belt:0.7.0 AS dev

ARG USER_NAME
ARG USER_UID
ARG USER_GID
ARG USER_GROUP
ARG REPO

COPY .docker/start.sh /usr/local/bin/start.sh

WORKDIR "/home/${USER_NAME}/src/${REPO}"

ENTRYPOINT [ "sh" ]

CMD [ "/usr/local/bin/start.sh" ]
