# <img src="https://i.imgur.com/z8BLePO.png" width="75%">

Pull based, language agnostic exception aggregator for microservice environments.

Periskop scales well with the number of exceptions and application instances:

  - Exceptions are pre-aggregated in client libraries and stored efficiently in memory, while keeping a sample of concrete occurrences for inspection.
  - Exceptions are scraped and aggregated across instances by the server component.
  - More application instances result in longer refresh cycles but the memory usage remains constant.

A UI component is provided for convenience.

See our talk "Periskop: Exception Monitoring at Scale" at FOSDEM [here](https://fosdem.org/2022/schedule/event/periskop/).

Related blog posts at SoundCloud's developers blog:

* [Periskop: Exception Monitoring Service](https://developers.soundcloud.com/blog/periskop-exception-monitoring-service)
* [What Is New with Periskop in 2022](https://developers.soundcloud.com/blog/periskop-in-2022)


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

## Pushgateway

See [periskop-pushgateway](https://github.com/periskop-dev/periskop-pushgateway) if you want to use Periskop as push based metric system.

## Client Libraries

  - [periskop-scala](https://github.com/periskop-dev/periskop-scala)
  - [periskop-go](https://github.com/periskop-dev/periskop-go)
  - [periskop-python](https://github.com/periskop-dev/periskop-python)
  - [periskop-ruby](https://github.com/periskop-dev/periskop-ruby)


## Integrations

  - [Backstage plugin](https://github.com/backstage/backstage/blob/master/plugins/periskop)

![backstage-plugin](https://github.com/backstage/backstage/blob/3131784fb71ac35ec24cea433a37bc5cab0595f9/plugins/periskop/docs/periskop-plugin-screenshot.png?raw=true)
