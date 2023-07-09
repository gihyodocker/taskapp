FROM golang:1.20.5

WORKDIR /go/src/github.com/gihyodocker/taskapp

COPY ./cmd ./cmd
COPY ./pkg ./pkg
COPY go.mod .
COPY go.sum .
COPY Makefile .
COPY ./assets ./assets

RUN make mod
RUN make vendor
RUN make build-web

ENTRYPOINT ["./bin/web"]