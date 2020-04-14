# <img src="https://i.imgur.com/NcT3qgf.png" width="65%">

[![Build Status](https://api.cirrus-ci.com/github/soundcloud/periskop.svg)](https://cirrus-ci.com/github/soundcloud/periskop)

Pull based, language agnostic exception aggregator for microservice environments.

Periskop scales well with the number of exceptions and application instances:

  - Exceptions are pre-aggregated in client libraries and stored efficiently in memory, while keeping a sample of concrete occurrences for inspection.
  - Exceptions are scraped and aggregated across instances by the server component.
  - More application instances result in longer refresh cycles but the memory usage remains constant.

A UI component is provided for convenience.

## Scraping

Errors are scraped and aggregated using a configured endpoint from each of the instances discovered via service discovery.

Periskop supports all service discovery mechanisms supported by Prometheus. The configuration format for service discovery
mirrors the one from Prometheus. See [Prometheus's official documentation](https://prometheus.io/docs/prometheus/latest/configuration/configuration/)
for reference.

At the moment, only DNS has been fully tested, we are working on properly testing the remaining methods.

A full example of service configuration for Periskop can be found in the [sample configuration](config.dev.yaml).

## Format

The format for scraped errors is defined in [a proto3 IDL](representation/errors.proto). Currently the only supported protocol is snake_cased JSON over HTTP ([example](scraper/sample-response1.json)).

## UI

The UI allows navigating and inspecting exceptions as they occur.

![ui](https://i.imgur.com/Tljxd80.png)

## Building & Running
We are looking into distributing Periskop via Docker Hub.
In the meantime, you can build and run Periskop from source:

```
docker build --tag periskop .
docker run -v path/to/config.yaml:/etc/periskop/periskop.yaml -p 8080:8080 periskop
```

## Client Libraries

  - [periskop-scala](https://github.com/soundcloud/periskop-scala)
  - [periskop-go](https://github.com/soundcloud/periskop-go)

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md)
