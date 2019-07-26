package repository

import (
	"fmt"
)

type ErrorAggregate struct {
	AggregationKey string             `json:"aggregation_key"`
	TotalCount     int                `json:"total_count"`
	LatestErrors   []ErrorWithContext `json:"latest_errors"`
}

type ErrorWithContext struct {
	Error       ErrorInstance `json:"error"`
	UUID        string        `json:"uuid"`
	Timestamp   int64         `json:"timestamp"`
	HTTPContext HTTPContext   `json:"http_context"`
}

type ErrorInstance struct {
	Class      string   `json:"class"`
	Message    string   `json:"message"`
	Stacktrace []string `json:"stacktrace"`
	Cause      *error   `json:"cause"`
}

type HTTPContext struct {
	RequestMethod  string            `json:"request_method"`
	RequestURL     string            `json:"request_url"`
	RequestHeaders map[string]string `json:"request_headers"`
}

type ErrorsRepository interface {
	StoreErrors(serviceName string, errors []ErrorAggregate)
	GetServices() []string
	GetErrors(serviceName string, numberOfErrors int) ([]ErrorAggregate, error)
}

func NewInMemory() ErrorsRepository {
	var r inMemoryRepository
	r.AggregatedError = make(map[string][]ErrorAggregate)
	return r
}

type inMemoryRepository struct {
	// map service name -> errors
	AggregatedError map[string][]ErrorAggregate
}

func (r inMemoryRepository) StoreErrors(serviceName string, errors []ErrorAggregate) {
	r.AggregatedError[serviceName] = errors
}

func (r inMemoryRepository) GetServices() []string {
	keys := make([]string, 0, len(r.AggregatedError))
	for k := range r.AggregatedError {
		keys = append(keys, k)
	}
	return keys
}

func (r inMemoryRepository) GetErrors(serviceName string, numberOfErrors int) ([]ErrorAggregate, error) {
	if value, ok := r.AggregatedError[serviceName]; ok {
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
	return nil, fmt.Errorf("service %s not found", serviceName)
}
