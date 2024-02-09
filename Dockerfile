ARG GO_VERSION=1.18
ARG ALPINE_VERSION=3.16

FROM golang:${GO_VERSION}-alpine as builder
RUN apk add --no-cache make gcc musl-dev git

WORKDIR $GOPATH/src/github.com/revolutionchain/charon
COPY go.mod go.sum $GOPATH/src/github.com/revolutionchain/charon/

# Cache go modules
RUN go mod download -x

ARG GIT_SHA
ENV CGO_ENABLED=0

COPY ./ $GOPATH/src/github.com/revolutionchain/charon

ENV GIT_SHA=$GIT_SH

RUN go build \
        -ldflags \
            "-X 'github.com/revolutionchain/charon/pkg/params.GitSha=`./sha.sh`'" \
        -o $GOPATH/bin $GOPATH/src/github.com/revolutionchain/charon/... && \
    rm -fr $GOPATH/src/github.com/revolutionchain/charon/.git

# Final stage
FROM alpine:${ALPINE_VERSION} as base
# Makefile supports generating ssl files from docker
RUN apk add --no-cache openssl
COPY --from=builder /go/bin/charon /charon

ENV REVO_RPC=http://revo:testpasswd@localhost:3889
ENV REVO_NETWORK=auto

EXPOSE 23889
EXPOSE 23890

ENTRYPOINT [ "/charon" ]