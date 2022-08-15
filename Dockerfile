ARG BUILDPLATFORM
ARG TARGETPLATFORM
ARG GO_VERSION=1.18
ARG ALPINE_VERSION=3.16

# Cross-compilation support
# FROM tonistiigi/xx AS xx

FROM golang:${GO_VERSION}-alpine as builder
COPY --from=xx / /
RUN apk add --no-cache make gcc musl-dev git

WORKDIR $GOPATH/src/github.com/qtumproject/janus
COPY go.mod go.sum $GOPATH/src/github.com/qtumproject/janus/

# Cache go modules
RUN go mod download -x

ARG GIT_SHA
ARG TARGETOS
ARG TARGETARCH
ARG TARGETPLATFORM
ENV CGO_ENABLED=0

COPY ./ $GOPATH/src/github.com/qtumproject/janus

ENV GIT_SHA=$GIT_SH

RUN build \
        -ldflags \
            "-X 'github.com/qtumproject/janus/pkg/params.GitSha=`./sha.sh`'" \
        -o $GOPATH/bin $GOPATH/src/github.com/qtumproject/janus/... && \
    rm -fr $GOPATH/src/github.com/qtumproject/janus/.git

# Final stage
FROM alpine:${ALPINE_VERSION} as base
# COPY --from=xx /out/xx-apk /out/xx-apk
COPY --from=xx /usr/bin/xx-apk /usr/bin/xx-apk
COPY --from=xx /usr/bin/xx-info /usr/bin/xx-info
# Makefile supports generating ssl files from docker
# RUN /usr/bin/xx-apk add --no-cache openssl
# RUN rm /usr/bin/xx-apk && rm /usr/bin/xx-info
COPY --from=builder /go/bin/janus /janus

ENV QTUM_RPC=http://qtum:testpasswd@localhost:3889
ENV QTUM_NETWORK=auto

EXPOSE 23889
EXPOSE 23890

ENTRYPOINT [ "/janus" ]