package repository

import (
	"testing"
)

func TestAddToResolved(t *testing.T) {
	serviceName := "test-service"
	er := &inMemoryRepository{}

	// service has empty list of resolved errors
	er.addToResolved(serviceName, "test-error-0")
	if value, ok := er.ResolvedErrors.Load(serviceName); ok {
		value, _ := value.([]string)
		if len(value) != 1 {
			t.Errorf("Expected 1 element, Found %d", len(value))
		}
	}

	// service has a resolved error
	er.addToResolved(serviceName, "test-error-1")
	if value, ok := er.ResolvedErrors.Load(serviceName); ok {
		value, _ := value.([]string)
		if len(value) != 2 {
			t.Errorf("Expected 2 element, Found %d", len(value))
		}
	}
}

func TestSearchResolved(t *testing.T) {
	serviceName := "test-service"
	er := &inMemoryRepository{}
	er.ResolvedErrors.Store(serviceName, []string{"test-error-0"})

	if !er.SearchResolved(serviceName, "test-error-0") {
		t.Errorf("Error should be found in resolved errors")
	}

	if er.SearchResolved(serviceName, "test-error-1") {
		t.Errorf("Error shouldn't be found in resolved errors")
	}
}

func TestResolveErrorError(t *testing.T) {
	serviceName := "test-service"
	er := &inMemoryRepository{}
	er.AggregatedError.Store(serviceName, []ErrorAggregate{
		{AggregationKey: "test-error-0"},
		{AggregationKey: "test-error-1"},
	})
	err := er.ResolveError(serviceName, "test-error-0")
	if err != nil {
		t.Errorf("deleting the error")
	}
	if value, ok := er.AggregatedError.Load(serviceName); ok {
		value, _ := value.([]ErrorAggregate)
		if len(value) != 1 {
			t.Errorf("Expected 1 element, Found %d", len(value))
		}
	}
}
