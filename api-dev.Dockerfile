## Build api
FROM golang:1.17.7

WORKDIR /periskop-dev

RUN mkdir /periskop-be

COPY go.mod .
COPY go.sum .
RUN go mod download

ENV PORT 8080
ENV CONFIG_FILE ./config.docker.yaml
