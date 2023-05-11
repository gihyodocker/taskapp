FROM golang:1.20.4

WORKDIR /go/src/github.com/gihyodocker/taskapp
COPY . .

RUN make mod
RUN make build

ENTRYPOINT ["./bin/main", "backend"]