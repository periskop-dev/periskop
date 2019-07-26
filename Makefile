.PHONY: clean run-web

WEB_FOLDER := web
PORT := 7777

clean:
	rm -rf $(WEB_FOLDER)/dist $(WEB_FOLDER)/node_modules

setup-web: web/package.json web/package-lock.json
	npm install --prefix $(WEB_FOLDER)

ci-setup-web: web/package.json web/package-lock.json
	npm ci --prefix $(WEB_FOLDER)

build-web:
	npm run build:dist --prefix $(WEB_FOLDER)

run-api:
	GO111MODULE=on go build -o periskop && ./periskop -port=$(PORT) -config ./config.dev.yaml

run-web:
	npm start --prefix $(WEB_FOLDER)

run: build-web run-api

build-api:
	GO111MODULE=on go build ./...

test-api:
	GO111MODULE=on go test ./...

lint:
	golangci-lint run