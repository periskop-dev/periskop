package repository

import (
	"testing"
)

const serviceName = "test-service"

func TestAddToResolved(t *testing.T) {
	er := &inMemoryRepository{}

	// service has empty list of resolved errors
	er.addToResolved(serviceName, "test-error-0")
	if value, ok := er.ResolvedErrors.Load(serviceName); ok {
		value, _ := value.(map[string]bool)
		if len(value) != 1 {
			t.Errorf("Expected 1 element, Found %d", len(value))
		}
	}

	// service has a resolved error
	er.addToResolved(serviceName, "test-error-1")
	if value, ok := er.ResolvedErrors.Load(serviceName); ok {
		value, _ := value.(map[string]bool)
		if len(value) != 2 {
			t.Errorf("Expected 2 element, Found %d", len(value))
		}
	}
}

func TestRemoveResolved(t *testing.T) {
	er := &inMemoryRepository{}

	er.ResolvedErrors.Store(serviceName, map[string]bool{"test-error-0": true})
	er.RemoveResolved(serviceName, "test-error-0")

	if value, ok := er.ResolvedErrors.Load(serviceName); ok {
		value, _ := value.(map[string]bool)
		if len(value) != 0 {
			t.Errorf("Expected 0 element, Found %d", len(value))
		}
	}
}

func TestSearchResolved(t *testing.T) {
	er := &inMemoryRepository{}
	er.ResolvedErrors.Store(serviceName, map[string]bool{"test-error-0": true})

	if !er.SearchResolved(serviceName, "test-error-0") {
		t.Errorf("Error should be found in resolved errors")
	}

	if er.SearchResolved(serviceName, "test-error-1") {
		t.Errorf("Error shouldn't be found in resolved errors")
	}
}

func TestResolveErrorError(t *testing.T) {
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
