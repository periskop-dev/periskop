## Build api
FROM golang:1.16 AS api-builder

WORKDIR /periskop

COPY go.mod .
COPY go.sum .
RUN /usr/local/go/bin/go mod download

COPY . .
RUN /usr/local/go/bin/go build -o app .

## Build web
FROM node:lts AS web-builder

WORKDIR /periskop
COPY . .
RUN npm ci --prefix web
RUN npm run build:dist --prefix web

## Build final container
FROM gcr.io/distroless/base

ENV SERVER_URL localhost
ENV SERVER_PORT 8080
ENV PORT 8080
ENV CONFIG_FILE /etc/periskop/periskop.yaml

COPY --from=web-builder /periskop/web/dist /periskop/web/dist
COPY --from=api-builder /periskop/app /periskop/app

CMD ["/periskop/app"]
