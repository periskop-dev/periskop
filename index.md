# <img src="https://i.imgur.com/z8BLePO.png width="75%">

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

A full example of service configuration for Periskop can be found in the [sample configuration](config.dev.yaml).

## Format

The format for scraped errors is defined in [a proto3 IDL](representation/errors.proto). Currently the only supported protocol is snake_cased JSON over HTTP ([example](scraper/sample-response1.json)).

## UI

The UI allows navigating and inspecting exceptions as they occur.

![ui](https://i.imgur.com/Tljxd80.png)

## Client Libraries

  - [periskop-scala](https://github.com/soundcloud/periskop-scala)
  - [periskop-go](https://github.com/soundcloud/periskop-go)
  - [periskop-python](https://github.com/soundcloud/periskop-python)
