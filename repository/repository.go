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
}

func NewInMemory() ErrorsRepository {
	return &inMemoryRepository{
		AggregatedError: sync.Map{},
	}
}

type inMemoryRepository struct {
	// map service name -> errors
	AggregatedError sync.Map
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

func (r *inMemoryRepository) DeleteError(serviceName string, key string) {
	if value, ok := r.AggregatedError.Load(serviceName); ok {
		value, _ := value.([]ErrorAggregate)
		newValues := []ErrorAggregate{}
		for _, errorObj := range value {
			if errorObj.AggregationKey != key {
				newValues = append(newValues, errorObj)
			}
		}
		r.StoreErrors(serviceName, newValues)
	}
}
