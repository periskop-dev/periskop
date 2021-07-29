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

// map error key -> errorAggregate
type errorAggregateMap map[string]errorAggregate

// map target -> error key -> error total occurrences
type targetErrorsCountMap map[string]map[string]int

// map error key -> list of errorWithContext (latest errors)
type errorInstancesAccumulatorMap map[string][]errorWithContext

type Scraper struct {
	Resolver      servicediscovery.Resolver
	Repository    *repository.ErrorsRepository
	ServiceConfig config.Service
	processor     Processor
}

// NewScraper create a new scraper for a given service name
func NewScraper(resolver servicediscovery.Resolver, r *repository.ErrorsRepository,
	serviceConfig config.Service, processor Processor) Scraper {
	return Scraper{
		Resolver:      resolver,
		Repository:    r,
		ServiceConfig: serviceConfig,
		processor:     processor,
	}
}

func (errorAggregates errorAggregateMap) combine(serviceName string, r *repository.ErrorsRepository,
	rp responsePayload, targetErrorsCount targetErrorsCountMap, errorInstancesAccumulator errorInstancesAccumulatorMap) {
	for _, item := range rp.ErrorAggregate {
		if _, exists := targetErrorsCount[rp.Target]; !exists {
			targetErrorsCount[rp.Target] = make(map[string]int)
		}
		prevErrorInstances := errorInstancesAccumulator[item.AggregationKey]
		var errorCountDelta int
		lastestErrors := combineLastErrors(prevErrorInstances, item.LatestErrors)

		if existing, exists := errorAggregates[item.AggregationKey]; exists {
			prevCount := targetErrorsCount[rp.Target][item.AggregationKey]
			if prevCount <= item.TotalCount {
				// Set the CreatedAt of its oldest occurrence
				createdAt := existing.CreatedAt
				if item.CreatedAt.Before(createdAt) {
					createdAt = item.CreatedAt
				}

				errorCountDelta = item.TotalCount - prevCount
				errorAggregates[item.AggregationKey] = errorAggregate{
					TotalCount:     existing.TotalCount + errorCountDelta,
					AggregationKey: existing.AggregationKey,
					Severity:       item.Severity,
					LatestErrors:   lastestErrors,
					CreatedAt:      createdAt,
				}
				updateValues(item, errorCountDelta, lastestErrors,
					serviceName, r, rp,
					targetErrorsCount, errorInstancesAccumulator)
			} else {
				log.Printf("warning: count of errors for '%s' target is inconsistent: prev %d, current %d.",
					rp.Target,
					prevCount,
					item.TotalCount,
				)
				log.Printf(" Counters won't be updated\n")
			}
		} else {
			errorAggregates[item.AggregationKey] = item
			updateValues(item, item.TotalCount, lastestErrors,
				serviceName, r, rp,
				targetErrorsCount, errorInstancesAccumulator)
		}
	}
}

func updateValues(item errorAggregate, errorCountDelta int, latestErrors []errorWithContext,
	serviceName string, r *repository.ErrorsRepository, rp responsePayload,
	targetErrorsCount targetErrorsCountMap, errorInstancesAccumulator errorInstancesAccumulatorMap) {
	metrics.ErrorOccurrences.WithLabelValues(serviceName, item.Severity, rp.Target,
		item.AggregationKey).Add(float64(errorCountDelta))
	targetErrorsCount[rp.Target][item.AggregationKey] = item.TotalCount
	errorInstancesAccumulator[item.AggregationKey] = latestErrors
	// If an error that was previously mark as resolved is scrapped again
	// it's going to be added to list of errors
	(*r).RemoveResolved(serviceName, item.AggregationKey)
}

func combineLastErrors(first []errorWithContext, second []errorWithContext) []errorWithContext {
	combined := append(first, second...)
	sort.Sort(errorOccurrences(combined))
	return combined
}

// Scrape runs go routines scrapping the list of targets of this service,
// processes the errors and stores them into the repository.
func (scraper Scraper) Scrape() {
	serviceConfig := scraper.ServiceConfig
	resolutions := scraper.Resolver.Resolve()
	var resolvedAddresses = servicediscovery.EmptyResolvedAddresses()
	timer := time.NewTimer(scraper.ServiceConfig.Scraper.RefreshInterval)

	var targetErrorsCount = make(targetErrorsCountMap)
	var errorAggregates = make(errorAggregateMap)
	for {
		select {
		case newResult := <-resolutions:
			resolvedAddresses = newResult
			storeTargets(serviceConfig.Name, scraper.Repository, resolvedAddresses)
			log.Printf("Received new dns resolution result for %s. Address resolved: %d\n", serviceConfig.Name,
				len(resolvedAddresses.Addresses))

		case <-timer.C:
			timer.Stop()
			errorInstancesAccumulator := make(errorInstancesAccumulatorMap)
			for responsePayload := range scrapeInstances(resolvedAddresses.Addresses, serviceConfig.Scraper.Endpoint,
				scraper.processor) {
				errorAggregates.combine(serviceConfig.Name, scraper.Repository,
					responsePayload, targetErrorsCount, errorInstancesAccumulator)
			}
			storeErrors(serviceConfig.Name, scraper.Repository, errorAggregates)

			numInstances := len(resolvedAddresses.Addresses)
			numErrors := len(errorAggregates)
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

func storeErrors(serviceName string, r *repository.ErrorsRepository, errorAggregates errorAggregateMap) {
	errors := make([]repository.ErrorAggregate, 0, len(errorAggregates))
	for _, value := range errorAggregates {
		if !(*r).SearchResolved(serviceName, value.AggregationKey) {
			severity := severityWithFallback(value.Severity)
			errors = append(errors, repository.ErrorAggregate{
				AggregationKey: value.AggregationKey,
				Severity:       severity,
				TotalCount:     value.TotalCount,
				LatestErrors:   toRepositoryErrorsWithContent(value.LatestErrors),
				CreatedAt:      value.CreatedAt.Unix(),
			})
		}
	}
	(*r).StoreErrors(serviceName, errors)
}

func storeTargets(serviceName string, r *repository.ErrorsRepository, addr servicediscovery.ResolvedAddresses) {
	targets := make([]repository.Target, 0, len(addr.Addresses))
	for _, value := range addr.Addresses {
		targets = append(targets, repository.Target{
			Endpoint: value,
		})
	}
	(*r).StoreTargets(serviceName, targets)
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
