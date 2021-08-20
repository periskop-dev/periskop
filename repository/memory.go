package repository

import (
	"fmt"
	"sync"

	"github.com/soundcloud/periskop/metrics"
)

func NewMemoryRepository() ErrorsRepository {
	return &memoryRepository{
		AggregatedError: sync.Map{},
		ResolvedErrors:  sync.Map{},
		Targets:         sync.Map{},
	}
}

type memoryRepository struct {
	// map service name -> list of errors
	AggregatedError sync.Map
	// map service name -> set of resolved errors
	ResolvedErrors sync.Map
	// map service name -> list of scraped targets
	Targets sync.Map
}

func (r *memoryRepository) GetErrors(serviceName string, numberOfErrors int) ([]ErrorAggregate, error) {
	if value, ok := r.AggregatedError.Load(serviceName); ok {
		value, _ := value.([]ErrorAggregate)
		result := make([]ErrorAggregate, 0, len(value))
		for _, errorObj := range value {
			topCap := len(errorObj.LatestErrors)
			if numberOfErrors < topCap {
				topCap = numberOfErrors
			}
			errorObj.LatestErrors = errorObj.LatestErrors[0:topCap]
			result = append(result, errorObj)
		}

		return result, nil
	}
	metrics.ServiceErrors.WithLabelValues("service_not_found").Inc()
	return nil, fmt.Errorf("service %s not found", serviceName)
}

func (r *memoryRepository) StoreErrors(serviceName string, errors []ErrorAggregate) {
	r.AggregatedError.Store(serviceName, errors)
}

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
// TODO: rename variables
func (r *memoryRepository) ResolveError(serviceName string, key string) error {
	if value, ok := r.AggregatedError.Load(serviceName); ok {
		value, _ := value.([]ErrorAggregate)
		newValues := []ErrorAggregate{}
		for _, errorObj := range value {
			if errorObj.AggregationKey != key {
				newValues = append(newValues, errorObj)
			}
		}
		r.StoreErrors(serviceName, newValues)
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

func (r *memoryRepository) StoreTargets(serviceName string, targets []Target) {
	r.Targets.Store(serviceName, targets)
}

func (r *memoryRepository) GetTargets() map[string][]Target {
	targets := make(map[string][]Target)
	r.Targets.Range(func(key, value interface{}) bool {
		targets[key.(string)] = value.([]Target)
		return true
	})
	return targets
}
