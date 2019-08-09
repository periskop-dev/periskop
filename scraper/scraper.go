package scraper

import (
	"log"
	"sort"
	"sync"
	"time"

	"github.com/soundcloud/periskop/config"
	"github.com/soundcloud/periskop/repository"
	"github.com/soundcloud/periskop/servicediscovery"
)

type errorAggregateMap map[string]errorAggregate

func (errorAggregateMap errorAggregateMap) combine(aggregatedErrors []errorAggregate) {
	for _, item := range aggregatedErrors {
		if existing, exists := errorAggregateMap[item.AggregationKey]; exists {
			errorAggregateMap[item.AggregationKey] = errorAggregate{
				TotalCount:     existing.TotalCount + item.TotalCount,
				AggregationKey: existing.AggregationKey,
				LatestErrors:   combine(existing.LatestErrors, item.LatestErrors),
			}
		} else {
			errorAggregateMap[item.AggregationKey] = item
		}
	}
}

func combine(first []errorWithContext, second []errorWithContext) []errorWithContext {
	combined := append(first, second...)
	sort.Sort(errorOccurrences(combined))
	return combined
}

type Scraper struct {
	Resolver      servicediscovery.SRVResolver
	Repository    *repository.ErrorsRepository
	ServiceConfig config.Service
	processor     Processor
}

func NewScraper(resolver servicediscovery.SRVResolver, r *repository.ErrorsRepository,
	serviceConfig config.Service, processor Processor) Scraper {
	return Scraper{
		Resolver:      resolver,
		Repository:    r,
		ServiceConfig: serviceConfig,
		processor:     processor,
	}
}

// Scrape stuff
func (scraper Scraper) Scrape() {
	serviceConfig := scraper.ServiceConfig
	resolutions := scraper.Resolver.Resolve()
	var resolvedAddresses = servicediscovery.EmptyResolvedAddresses()
	timer := time.NewTimer(scraper.ServiceConfig.Scraper.RefreshInterval)

	for {
		select {

		case newResult := <-resolutions:
			resolvedAddresses = newResult
			log.Printf("Received new dns resolution result for %s. Address resolved: %d\n", serviceConfig.Name,
				len(resolvedAddresses.Addresses))

		case <-timer.C:
			timer.Stop()
			var currentAggregatedErrorsMap = make(errorAggregateMap)
			for responsePayload := range scrapeInstances(resolvedAddresses.Addresses, serviceConfig.Scraper.Endpoint,
				scraper.processor) {
				currentAggregatedErrorsMap.combine(responsePayload)
			}
			store(serviceConfig.Name, scraper.Repository, currentAggregatedErrorsMap)
			log.Printf("%s: scraped %d errors from %d instances", serviceConfig.Name,
				len(currentAggregatedErrorsMap), len(resolvedAddresses.Addresses))
			timer.Reset(scraper.ServiceConfig.Scraper.RefreshInterval)
		}
	}
}

func scrapeInstances(addresses []string, endpoint string, processor Processor) <-chan []errorAggregate {
	var wg sync.WaitGroup
	out := make(chan []errorAggregate, len(addresses))

	wg.Add(len(addresses))
	for _, address := range addresses {
		request := Request{
			Target:        "http://" + address + endpoint,
			ResultChannel: out,
			WaitGroup:     &wg,
		}

		go processor.Enqueue(request)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func store(serviceName string, r *repository.ErrorsRepository, m errorAggregateMap) {
	errors := make([]repository.ErrorAggregate, 0, len(m))
	for _, value := range m {
		errors = append(errors, repository.ErrorAggregate{
			AggregationKey: value.AggregationKey,
			TotalCount:     value.TotalCount,
			LatestErrors:   toRepositoryErrorsWithContent(value.LatestErrors),
		})
	}
	(*r).StoreErrors(serviceName, errors)
}

func toRepositoryErrorsWithContent(occurrences []errorWithContext) []repository.ErrorWithContext {
	errors := make([]repository.ErrorWithContext, 0, len(occurrences))
	for _, occurrence := range occurrences {
		errors = append(errors, repository.ErrorWithContext{
			Timestamp: occurrence.Timestamp.Unix(),
			UUID:      occurrence.UUID,
			Error: repository.ErrorInstance{
				Class:      occurrence.Error.Class,
				Message:    occurrence.Error.Message,
				Stacktrace: occurrence.Error.Stacktrace,
				Cause:      nil, //TODO map the cause recursively
			},
			HTTPContext: repository.HTTPContext{
				RequestHeaders: occurrence.HTTPContext.RequestHeaders,
				RequestMethod:  occurrence.HTTPContext.RequestMethod,
				RequestURL:     occurrence.HTTPContext.RequestURL,
			},
		})
	}
	return errors
}
