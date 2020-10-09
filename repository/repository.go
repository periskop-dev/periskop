package repository

import (
	"fmt"
	"sync"

	"github.com/soundcloud/periskop/metrics"
)

type ErrorAggregate struct {
	AggregationKey string             `json:"aggregation_key"`
	TotalCount     int                `json:"total_count"`
	Severity       string             `json:"severity"`
	LatestErrors   []ErrorWithContext `json:"latest_errors"`
	CreatedAt      int64              `json:"created_at"`
}

type ErrorWithContext struct {
	Error       ErrorInstance `json:"error"`
	UUID        string        `json:"uuid"`
	Timestamp   int64         `json:"timestamp"`
	Severity    string        `json:"severity"`
	HTTPContext *HTTPContext  `json:"http_context"`
}

type ErrorInstance struct {
	Class      string         `json:"class"`
	Message    string         `json:"message"`
	Stacktrace []string       `json:"stacktrace"`
	Cause      *ErrorInstance `json:"cause"`
}

type HTTPContext struct {
	RequestMethod  string            `json:"request_method"`
	RequestURL     string            `json:"request_url"`
	RequestHeaders map[string]string `json:"request_headers"`
	RequestBody    string            `json:"request_body"`
}

type ErrorsRepository interface {
	StoreErrors(serviceName string, errors []ErrorAggregate)
	GetServices() []string
	GetErrors(serviceName string, numberOfErrors int) ([]ErrorAggregate, error)
	ResolveError(serviceName string, key string) error
	SearchResolved(serviceName string, key string) bool
	RemoveResolved(serviceName string, key string)
}

func NewInMemory() ErrorsRepository {
	return &inMemoryRepository{
		AggregatedError: sync.Map{},
		ResolvedErrors:  sync.Map{},
	}
}

type inMemoryRepository struct {
	// map service name -> list of errors
	AggregatedError sync.Map
	// map service name -> set of resolved errors
	ResolvedErrors sync.Map
}

func (r *inMemoryRepository) StoreErrors(serviceName string, errors []ErrorAggregate) {
	r.AggregatedError.Store(serviceName, errors)
}

func (r *inMemoryRepository) GetServices() []string {
	keys := make([]string, 0)
	r.AggregatedError.Range(func(key, value interface{}) bool {
		k, _ := key.(string)
		keys = append(keys, k)
		return true
	})
	return keys
}

func (r *inMemoryRepository) GetErrors(serviceName string, numberOfErrors int) ([]ErrorAggregate, error) {
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

// addToResolved saves the error to resolved error set
func (r *inMemoryRepository) addToResolved(serviceName string, key string) {
	if resolvedSet, ok := r.ResolvedErrors.Load(serviceName); ok {
		resolvedSet := resolvedSet.(map[string]bool)
		resolvedSet[key] = true
		r.ResolvedErrors.Store(serviceName, resolvedSet)
	} else {
		r.ResolvedErrors.Store(serviceName, map[string]bool{key: true})
	}
}

// RemoveResolved removes a resolved error from resolved error set
func (r *inMemoryRepository) RemoveResolved(serviceName string, key string) {
	if resolvedSet, ok := r.ResolvedErrors.Load(serviceName); ok {
		resolvedSet := resolvedSet.(map[string]bool)
		delete(resolvedSet, key)
		r.ResolvedErrors.Store(serviceName, resolvedSet)
	}
}

// SearchResolved searches if an error is inside the set of resolved errors
func (r *inMemoryRepository) SearchResolved(serviceName string, key string) bool {
	if resolvedSet, ok := r.ResolvedErrors.Load(serviceName); ok {
		resolvedSet := resolvedSet.(map[string]bool)
		return resolvedSet[key]
	}
	return false
}

// ResolveError removes the error from list of errors and adds to the set of resolved errors
func (r *inMemoryRepository) ResolveError(serviceName string, key string) error {
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
