## Build api
FROM golang:1.17.7

WORKDIR /mock-target

COPY . .
RUN go build -o mock-target mocktarget.go

ENTRYPOINT [ "./mock-target" ]
