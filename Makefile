.PHONY: clean setup-web ci-setup-web build-web run-api run-mock-target run-web run build-api test-api lint-api

WEB_FOLDER := web
PORT := 7777
DOCKER_IP := $(shell docker-machine ip default)

clean:
	rm -rf $(WEB_FOLDER)/dist $(WEB_FOLDER)/node_modules

setup-web: web/package.json web/package-lock.json
	npm install --prefix $(WEB_FOLDER)

ci-setup-web: web/package.json web/package-lock.json
	npm ci --prefix $(WEB_FOLDER)

build-web:
	npm run build:dist --prefix $(WEB_FOLDER)

run-api:
	GO111MODULE=on go build -o periskop && SERVER_URL=localhost ./periskop -port=$(PORT) -config ./config.dev.yaml

run-mock-target:
	cd mocktarget && GO111MODULE=on go build -o mock-target mocktarget.go && ./mock-target

run-web:
	npm start --prefix $(WEB_FOLDER)

run: build-web run-api

build-api:
	GO111MODULE=on go build ./...

test-api:
	GO111MODULE=on go test ./...

lint-api:
	golangci-lint run

###########################################################################################
## DOCKER COMPOSE
###########################################################################################

down:
	docker-compose down

build-up:
	SERVER_URL=$(DOCKER_IP) docker-compose up --build -d	

up: clean
	SERVER_URL=$(DOCKER_IP) docker-compose up -d

logs:
	docker-compose logs