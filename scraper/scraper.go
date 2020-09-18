package scraper

import (
	"log"
	"sort"
	"sync"
	"time"

	"github.com/soundcloud/periskop/config"
	"github.com/soundcloud/periskop/metrics"
	"github.com/soundcloud/periskop/repository"
	"github.com/soundcloud/periskop/servicediscovery"
)

type errorAggregateMap map[string]errorAggregate
type instanceErrorAggregateMap map[string]map[string]int

func (ea errorAggregateMap) combine(rp responsePayload, es instanceErrorAggregateMap) {
	for _, item := range rp.ErrorAggregate {
		if existing, exists := ea[item.AggregationKey]; exists {
			prevCount := es[rp.Instance][item.AggregationKey]

			ea[item.AggregationKey] = errorAggregate{
				TotalCount:     existing.TotalCount + (item.TotalCount - prevCount),
				AggregationKey: existing.AggregationKey,
				Severity:       item.Severity,
				LatestErrors:   combine(existing.LatestErrors, item.LatestErrors),
			}
			es[rp.Instance][item.AggregationKey] = item.TotalCount
		} else {
			ea[item.AggregationKey] = item
			es[rp.Instance] = make(map[string]int)
			es[rp.Instance][item.AggregationKey] = 0
		}
	}
}

func combine(first []errorWithContext, second []errorWithContext) []errorWithContext {
	combined := append(first, second...)
	sort.Sort(errorOccurrences(combined))
	return combined
}

type Scraper struct {
	Resolver      servicediscovery.Resolver
	Repository    *repository.ErrorsRepository
	ServiceConfig config.Service
	processor     Processor
}

func NewScraper(resolver servicediscovery.Resolver, r *repository.ErrorsRepository,
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
	var errorsSnapshot = make(instanceErrorAggregateMap)
	var currentAggregatedErrorsMap = make(errorAggregateMap)
	for {
		select {
		case newResult := <-resolutions:
			resolvedAddresses = newResult
			log.Printf("Received new dns resolution result for %s. Address resolved: %d\n", serviceConfig.Name,
				len(resolvedAddresses.Addresses))

		case <-timer.C:
			timer.Stop()
			for responsePayload := range scrapeInstances(resolvedAddresses.Addresses, serviceConfig.Scraper.Endpoint,
				scraper.processor) {
				currentAggregatedErrorsMap.combine(responsePayload, errorsSnapshot)
			}
			store(serviceConfig.Name, scraper.Repository, currentAggregatedErrorsMap)
			numInstances := len(resolvedAddresses.Addresses)
			numErrors := len(currentAggregatedErrorsMap)
			metrics.InstancesScrapped.WithLabelValues(serviceConfig.Name).Set(float64(numInstances))
			metrics.ErrorsScrapped.WithLabelValues(serviceConfig.Name).Add(float64(numErrors))
			log.Printf("%s: scraped %d errors from %d instances", serviceConfig.Name, numErrors, numInstances)
			timer.Reset(scraper.ServiceConfig.Scraper.RefreshInterval)
		}
	}
}

func scrapeInstances(addresses []string, endpoint string, processor Processor) <-chan responsePayload {
	var wg sync.WaitGroup
	out := make(chan responsePayload, len(addresses))

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
		severity := severityWithFallback(value.Severity)
		errors = append(errors, repository.ErrorAggregate{
			AggregationKey: value.AggregationKey,
			Severity:       severity,
			TotalCount:     value.TotalCount,
			LatestErrors:   toRepositoryErrorsWithContent(value.LatestErrors),
		})
		metrics.ErrorOccurrences.WithLabelValues(serviceName, severity,
			value.AggregationKey).Set(float64(value.TotalCount))
	}
	(*r).StoreErrors(serviceName, errors)
}

func toRepositoryErrorsWithContent(occurrences []errorWithContext) []repository.ErrorWithContext {
	errors := make([]repository.ErrorWithContext, 0, len(occurrences))
	for _, occurrence := range occurrences {
		errors = append(errors, repository.ErrorWithContext{
			Timestamp: occurrence.Timestamp.Unix(),
			Severity:  severityWithFallback(occurrence.Severity),
			UUID:      occurrence.UUID,
			Error: repository.ErrorInstance{
				Class:      occurrence.Error.Class,
				Message:    occurrence.Error.Message,
				Stacktrace: occurrence.Error.Stacktrace,
				Cause:      toRepositoryErrorCause(&occurrence.Error),
			},
			HTTPContext: toRepositoryHTTPContext(occurrence.HTTPContext),
		})
	}
	return errors
}

func severityWithFallback(severity string) string {
	if severity == "" {
		return "error"
	}
	return severity
}

func toRepositoryErrorCause(errorInstance *errorInstance) *repository.ErrorInstance {
	if errorInstance.Cause == nil {
		return nil
	}
	return &repository.ErrorInstance{
		Class:      errorInstance.Cause.Class,
		Message:    errorInstance.Cause.Message,
		Stacktrace: errorInstance.Cause.Stacktrace,
		Cause:      toRepositoryErrorCause(errorInstance.Cause),
	}
}

func toRepositoryHTTPContext(httpContext *httpContext) *repository.HTTPContext {
	if httpContext == nil {
		return nil
	}

	return &repository.HTTPContext{
		RequestHeaders: httpContext.RequestHeaders,
		RequestMethod:  httpContext.RequestMethod,
		RequestURL:     httpContext.RequestURL,
		RequestBody:    httpContext.RequestBody,
	}
}
