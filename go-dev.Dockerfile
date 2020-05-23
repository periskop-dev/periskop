## Build api
FROM golang:1.14.1

WORKDIR /periskop-dev

RUN mkdir /periskop-be

ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .
RUN go mod download

ENV PORT 8080
ENV SERVER_URL localhost
ENV CONFIG_FILE ./config.dev.yaml