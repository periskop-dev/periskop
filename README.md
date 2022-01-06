# <img src="https://i.imgur.com/z8BLePO.png" width="65%">

[![Build Status](https://api.cirrus-ci.com/github/soundcloud/periskop.svg)](https://cirrus-ci.com/github/soundcloud/periskop)
[![Gitter chat](https://badges.gitter.im/gitterHQ/gitter.png)](https://gitter.im/soundcloud/periskop)

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

## Run project locally

Please see [CONTRIBUTING.md](CONTRIBUTING.md)

## Building & Running

We are looking into distributing Periskop via Docker Hub.
In the meantime, you can build and run Periskop from source:

```
docker build --tag periskop .
docker run -v path/to/config.yaml:/etc/periskop/periskop.yaml -p 8080:8080 periskop
```

## Enable persistance storage

By default Periskop stores all the scrapped errors in memory [repository](repository/memory.go). You can configure your Periskop deployment to use persistent storage.
Currently the supported persistance storages are SQLite, MySQL and PostgreSQL. 

For SQLite, add these lines to your `config.yaml` file:
```yaml
repository:
  type: sqlite
  path: periskop.db
```

For MySQL:
```yaml
repository:
  type: mysql
  dsn: user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
```
  
For PostgreSQL:
```yaml
repository:
  type: postgres
  dsn: host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable
```

## Alert reported exceptions

All reported errors are instrumented with [Prometheus](https://prometheus.io) which provides alerting capabilities using [Alertmanager](https://prometheus.io/docs/alerting/alertmanager/). You can configure an alert when you reach some threshold of errors. Here's an example:

```yaml
groups:
- name: periskop
  rules:
  - alert: TooManyErrors
    expr: periskop_error_occurrences{severity="error"} > 1000
    for: 5m
    labels:
      severity: critical    
    annotations:
      summary: "Too many errors on {{ $labels.service_name }}"
      description: "Errors for {{ $labels.service_name }}({{ $labels.aggregation_key }}) is {{ $value }}"
      dashboard: "https://periskop.example.com/#/{{ $labels.service_name }}/errors/{{ $labels.aggregation_key }}"
```

## Pushgateway

See [periskop-pushgateway](https://github.com/soundcloud/periskop-pushgateway) if you want to use Periskop as push based metric system.

## Client Libraries

  - [periskop-scala](https://github.com/periskop-dev/periskop-scala)
  - [periskop-go](https://github.com/periskop-dev/periskop-go)
  - [periskop-python](https://github.com/periskop-dev/periskop-python)
  - [periskop-ruby](https://github.com/periskop-dev/periskop-ruby)
