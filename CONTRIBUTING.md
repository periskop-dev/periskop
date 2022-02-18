# Contributing

## Setup

Periskop requires `go`, `make` and `npm` installed in your system. Once they are installed, you can install the remaining dependencies by running:

```bash
make setup-web
```

Make sure tests pass:

```bash
make test-api
```

## Development

Periskop can be run in development in two different ways, either by running the API and web components separately, or together. The difference is that
by running them separately the webpack dev server will be used, with support for live reloading.

Running the API:

```bash
make run-api
```

Running the web component:

```bash
make run-web
```

Periskop needs something to scrape in order to show errors in the UI. In development, a mock target returning
static responses can be used for this purpose (periskop is configured to point to it by default). To run the mock
server:

```bash
make run-mock-target
```

Now you can point your browser to `http://localhost:3000`

### Using `docker-compose`

We also provide a way to run periskop using `docker` container and orchestrated via composition.

In the `Makefile` there's a section for `DOCKER COMPOSE`:

- `make up`: Boots up the 3 containers (api, web, mock-target). Note: I may take some time until all the containers boot up and errors scraped and shown in the UI.
- `make down`: Will call `docker-compose down` which will stop all running containers orchestrated with our `docker-compose` configuration
- `make logs`: Shows the logs of the 3 containers

If you would like to follow the logs of any container you can do so as follows:

```bash
# Front end container - nodejs
docker-compose logs --follow web

# Back end container - golang
docker-compose logs --follow api

# The mock target. Doesn't do much other that show 1 line of code where it is serving the errors - golang
docker-compose logs --follow mock-target
```

All the commands that boot up containers override 2 environment variables:

- `API_HOST`: This is the host of the api. 
- `API_PORT`: The port where the api will be listening to. It defaults to `8080`.

## Testing

Running the API tests:

```bash
make test-api
```

Running the API linter:

```bash
make lint-api
```

## Contribution Policy

1. Fork the project and implement your changes in a separate branch. Make sure tests and linters pass.
1. Open a PR from your branch to periskop. Make sure the description is concise and clear.
1. Wait for an approval from 2 core developers. Make sure the build is green.
1. After that, one of the core developers will merge your changes into the main branch.
1. Thank you so much for contributing!
