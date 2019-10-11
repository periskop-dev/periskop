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

Run everything (note that a different port will be used):

```bash
make run
```

Periskop needs something to scrape in order to show errors in the UI. In development, a mock target returning
static responses can be used for this purpose (periskop is configured to point to it by default). To run the mock
server:

```bash
make run-mock-target
```

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
