FROM golang:1.20.5

WORKDIR /migrator

RUN apt update
RUN apt install -y default-mysql-client
RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.16.2

COPY . .