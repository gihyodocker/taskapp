FROM golang:1.21.6

WORKDIR /go/src/github.com/gihyodocker/taskapp

COPY ./cmd ./cmd
COPY ./pkg ./pkg
COPY go.mod .
COPY go.sum .
COPY Makefile .

RUN make mod
RUN make vendor
RUN make build-api

ENTRYPOINT ["./bin/api"]
