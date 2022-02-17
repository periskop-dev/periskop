package repository

import (
	"fmt"
	"sync"

	"github.com/periskop-dev/periskop/metrics"
)

func NewMemoryRepository() ErrorsRepository {
	return &memoryRepository{
		AggregatedError: sync.Map{},
		ResolvedErrors:  sync.Map{},
	}
}

type memoryRepository struct {
	// map service name -> list of errors
	AggregatedError sync.Map
	// map service name -> set of resolved errors
	ResolvedErrors sync.Map
	targetsRepository
}

// GetErrors fetches the last numberOfErrors of each aggregation of errors for the given service
func (r *memoryRepository) GetErrors(serviceName string, numberOfErrors int) ([]ErrorAggregate, error) {
	if value, ok := r.AggregatedError.Load(serviceName); ok {
		prevErrors, _ := value.([]ErrorAggregate)
		errors := make([]ErrorAggregate, 0, len(prevErrors))
		for _, errorAggregate := range prevErrors {
			maxErrors := len(errorAggregate.LatestErrors)
			if numberOfErrors < maxErrors {
				maxErrors = numberOfErrors
			}
			errorAggregate.LatestErrors = errorAggregate.LatestErrors[0:maxErrors]
			errors = append(errors, errorAggregate)
		}

		return errors, nil
	}
	metrics.ServiceErrors.WithLabelValues("service_not_found").Inc()
	return nil, fmt.Errorf("service %s not found", serviceName)
}

// ReplaceErrors replaces a lists of aggregated errors for the given service
func (r *memoryRepository) ReplaceErrors(serviceName string, errors []ErrorAggregate) {
	r.AggregatedError.Store(serviceName, errors)
}

// GetServices fetches the list of unique services
func (r *memoryRepository) GetServices() []string {
	keys := make([]string, 0)
	r.AggregatedError.Range(func(key, value interface{}) bool {
		k, _ := key.(string)
		keys = append(keys, k)
		return true
	})
	return keys
}

// ResolveError removes the error from list of errors and adds to the set of resolved errors
func (r *memoryRepository) ResolveError(serviceName string, key string) error {
	if value, ok := r.AggregatedError.Load(serviceName); ok {
		prevErrors, _ := value.([]ErrorAggregate)
		errors := []ErrorAggregate{}
		for _, errorAggregate := range prevErrors {
			if errorAggregate.AggregationKey != key {
				errors = append(errors, errorAggregate)
			}
		}
		r.ReplaceErrors(serviceName, errors)
		r.addToResolved(serviceName, key)
		return nil
	}
	return fmt.Errorf("service %s not found", serviceName)
}

// addToResolved saves the error to resolved error set
func (r *memoryRepository) addToResolved(serviceName string, key string) {
	if resolvedSet, ok := r.ResolvedErrors.Load(serviceName); ok {
		resolvedSet := resolvedSet.(map[string]bool)
		resolvedSet[key] = true
		r.ResolvedErrors.Store(serviceName, resolvedSet)
	} else {
		r.ResolvedErrors.Store(serviceName, map[string]bool{key: true})
	}
}

// RemoveResolved removes a resolved error from resolved error set
func (r *memoryRepository) RemoveResolved(serviceName string, key string) {
	if resolvedSet, ok := r.ResolvedErrors.Load(serviceName); ok {
		resolvedSet := resolvedSet.(map[string]bool)
		delete(resolvedSet, key)
		r.ResolvedErrors.Store(serviceName, resolvedSet)
	}
}

// SearchResolved searches if an error is inside the set of resolved errors
func (r *memoryRepository) SearchResolved(serviceName string, key string) bool {
	if resolvedSet, ok := r.ResolvedErrors.Load(serviceName); ok {
		resolvedSet := resolvedSet.(map[string]bool)
		return resolvedSet[key]
	}
	return false
}
