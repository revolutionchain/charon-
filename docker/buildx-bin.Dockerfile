FROM docker
COPY --from=docker/buildx-bin /buildx /usr/libexec/docker/cli-plugins/docker-buildx
ENTRYPOINT [ "docker" ]
