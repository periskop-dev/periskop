.PHONY: clean setup-web ci-setup-web build-web run-api run-mock-target run-web run build-api test-api lint-api

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
	go build -o periskop && ./periskop -port=$(PORT) -config ./config.dev.yaml

run-mock-target:
	go build -o mocktarget/mock-target mocktarget/mocktarget.go && cd mocktarget && ./mock-target

run-web:
	npm start --prefix $(WEB_FOLDER)

run: build-web run-api

build-api:
	go build ./...

test-api:
	go test ./...

lint-api:
	golangci-lint run

###########################################################################################
## DOCKER COMPOSE
###########################################################################################

down:
	docker-compose down

up: clean
	docker-compose up -d

logs:
	docker-compose logs -f
