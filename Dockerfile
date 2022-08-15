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

COPY ./ $GOPATH/src/github.com/qtumproject/janus

ENV GIT_SHA=$GIT_SH

RUN go build \
        -ldflags \
            "-X 'github.com/qtumproject/janus/pkg/params.GitSha=`./sha.sh`'" \
        -o $GOPATH/bin $GOPATH/src/github.com/qtumproject/janus/... && \
    rm -fr $GOPATH/src/github.com/qtumproject/janus/.git

# Final stage
FROM alpine:${ALPINE_VERSION} as base
# Makefile supports generating ssl files from docker
RUN apk add --no-cache openssl
COPY --from=builder /go/bin/janus /janus

ENV QTUM_RPC=http://qtum:testpasswd@localhost:3889
ENV QTUM_NETWORK=auto

EXPOSE 23889
EXPOSE 23890

ENTRYPOINT [ "/janus" ]