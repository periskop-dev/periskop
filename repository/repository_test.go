package repository

import (
	"testing"
)

func TestDeleteError(t *testing.T) {
	serviceName := "test-service"
	er := &inMemoryRepository{}
	er.AggregatedError.Store(serviceName, []ErrorAggregate{
		{AggregationKey: "test-error-0"},
		{AggregationKey: "test-error-1"},
	})
	err := er.DeleteError(serviceName, "test-error-0")
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
