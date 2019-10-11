# Periskop

[![Build Status](https://api.cirrus-ci.com/github/soundcloud/periskop.svg)](https://cirrus-ci.com/github/soundcloud/periskop)

Pull based, language agnostic exception aggregator for microservice environments.

Periskop scales well with the number of exceptions and application instances:

  - Exceptions are pre-aggregated in client libraries and stored efficiently in memory, while keeping a sample of concrete occurrences for inspection.
  - Exceptions are scraped and aggregated across instances by the server component.
  - More application instances result in longer refresh cycles but the memory usage remains constant.
  - Support for delegation is planned.

A UI component is provided for convenience.

## Scraping

Errors are scraped and aggregated using a configured endpoint from each of the instances discovered via service discovery.

At the moment, only DNS service discovery is supported. See the [sample configuration](config.dev.yaml).

## Format

The format for scraped errors is defined in [a proto3 IDL](representation/errors.proto). Currently the only supported protocol is snake_cased JSON over HTTP ([example](scraper/sample-response1.json)).

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md)
