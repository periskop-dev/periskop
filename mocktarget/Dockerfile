## Build api
FROM golang:1.14.1

WORKDIR /mock-target

ENV GO111MODULE=on

COPY . .
RUN go build -o mock-target mocktarget.go

ENTRYPOINT [ "./mock-target" ]