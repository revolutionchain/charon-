FROM golang:1.14-alpine

RUN echo $GOPATH
RUN apk add --no-cache make gcc musl-dev git
WORKDIR $GOPATH/src/github.com/qtumproject/truffle-parser
COPY ./main.go $GOPATH/src/github.com/qtumproject/truffle-parser
RUN go get -d ./...
RUN go install github.com/qtumproject/truffle-parser/

ENTRYPOINT [ "truffle-parser" ]