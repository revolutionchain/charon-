FROM golang:1.18

WORKDIR $GOPATH/src/github.com/revolutionchain/charon
COPY . $GOPATH/src/github.com/revolutionchain/charon
RUN go get -d ./...

CMD [ "go", "test", "-v", "./..."]