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

At the moment, only DNS service discovery is supported. See the [sample configuration](config.dev.yaml).

## Format

The format for scraped errors is defined in [a proto3 IDL](representation/errors.proto). Currently the only supported protocol is snake_cased JSON over HTTP ([example](scraper/sample-response1.json)).

## UI

The UI allows navigating and inspecting exceptions as they occur.

![ui](https://i.imgur.com/Tljxd80.png)

## Client Libraries

  - [periskop-scala](https://github.com/soundcloud/periskop-scala)
  - [periskop-go](https://github.com/soundcloud/periskop-go)

## Alert reported exceptions

All reported errors are instrumented with [Prometheus](https://prometheus.io) which provides alerting capabilities using [Alertmanager](https://prometheus.io/docs/alerting/alertmanager/). You can configure an alert when you reach some threshold of errors. Here's an example:

```yaml
groups:
- name: periskop
  rules:
  - alert: TooManyErrors
    expr: periskop_scrapped_aggregated_error_total{severity="error"} > 1000
    for: 5m
    labels:
      severity: critical    
    annotations:
      summary: "Too many errors on {{ $labels.service_name }}"
      description: "Errors for {{ $labels.service_name }}({{ $labels.aggregation_key }}) is {{ $value }}"
      dashboard: "https://periskop.example.com/#/{{ $labels.service_name }}/errors/{{ $labels.aggregation_key }}"
```

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md)
