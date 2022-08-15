ARG GO_VERSION=1.18
ARG ALPINE_VERSION=3.16

FROM golang:${GO_VERSION}-alpine as builder
RUN apk add --no-cache make gcc musl-dev git

WORKDIR $GOPATH/src/github.com/qtumproject/janus
COPY go.mod go.sum $GOPATH/src/github.com/qtumproject/janus/

# Cache go modules
RUN go mod download -x

ARG GIT_SHA
ENV CGO_ENABLED=0
ENV GIT_SHA=$GIT_SH

ENTRYPOINT [ "go" ]
