ARG BUILDPLATFORM=linux/amd64
ARG GO_VERSION=1.18
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS build

RUN apk add --no-cache git g++ openssh-client
RUN mkdir -p $GOPATH/src/github.com/qdm12/xcputranslate
RUN git clone https://github.com/qdm12/xcputranslate $GOPATH/src/github.com/qdm12/xcputranslate

WORKDIR $GOPATH/src/github.com/qdm12/xcputranslate

# 0.7.0
RUN git checkout ce3ad75f269b6c7097f26bd68348791d0b9a16e8
RUN go build -o /usr/local/bin/xcputranslate cmd/xcputranslate/main.go

ENTRYPOINT ["/usr/local/bin/xcputranslate"]
# USER 1000
# COPY --from=build --chown=1000 /tmp/gobuild/entrypoint /xcputranslate
# COPY --from=build --chown=1000 /usr/local/bin/xcputranslate /xcputranslate
