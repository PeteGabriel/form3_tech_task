FROM golang:1.16 as builder

ENV GO111MODULE=on
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

CMD go test -v ./...