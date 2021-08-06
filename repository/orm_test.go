package repository

import (
	"reflect"
	"testing"
	"time"
)

func TestStoreErrors(t *testing.T) {
	r := NewORMRepository()
	errors := []ErrorAggregate{
		{
			AggregationKey: "key",
			Severity:       "error",
			CreatedAt:      time.Unix(0, 0).Unix(),
			LatestErrors: []ErrorWithContext{
				{
					Error:     ErrorInstance{},
					Severity:  "error",
					Timestamp: time.Unix(0, 0).Unix(),
				},
			},
		}}
	r.StoreErrors("test", errors)
	r.StoreErrors("test", errors)
	if r.countErrors("test") != 1 {
		t.Errorf("Found %d errors, expected 1", r.countErrors("test"))
	}
}

func TestGetErrors(t *testing.T) {
	r := NewORMRepository()
	err0 := ErrorAggregate{
		AggregationKey: "key0",
		Severity:       "error",
		CreatedAt:      time.Unix(0, 0).Unix(),
		LatestErrors: []ErrorWithContext{
			{
				Error:     ErrorInstance{},
				Severity:  "error",
				Timestamp: time.Unix(0, 0).Unix(),
			},
		},
	}
	err1 := ErrorAggregate{
		AggregationKey: "key1",
		Severity:       "error",
		CreatedAt:      time.Unix(0, 0).Unix(),
		LatestErrors: []ErrorWithContext{
			{
				Error: ErrorInstance{
					Class:      "ErrorString",
					Message:    "Failed when parsing",
					Stacktrace: []string{"0:", "error"},
				},
				Severity:  "error",
				Timestamp: time.Unix(0, 0).Unix(),
			},
		},
	}
	errors := []ErrorAggregate{err0, err1}

	r.StoreErrors("test", errors)

	aggregatedErrors, err := r.GetErrors("test", 5)
	if err != nil {
		t.Errorf("Fail to fetch errors: %s", err)
	}

	if len(aggregatedErrors) != 2 {
		t.Errorf("Found %d errors, expected 2", len(aggregatedErrors))
	}

	if !reflect.DeepEqual(aggregatedErrors[0], err0) {
		t.Errorf("Error fetching errors, got %+v errors, expected %+v", &aggregatedErrors[0], &err0)
	}

	if !reflect.DeepEqual(aggregatedErrors[1], err1) {
		t.Errorf("Error fetching errors, got %+v errors, expected %+v", &aggregatedErrors[1], &err1)
	}
}

func TestGetServices(t *testing.T) {
	r := NewORMRepository()
	errors := []ErrorAggregate{
		{
			AggregationKey: "key",
			Severity:       "error",
			CreatedAt:      time.Unix(0, 0).Unix(),
			LatestErrors: []ErrorWithContext{
				{
					Error:     ErrorInstance{},
					Severity:  "error",
					Timestamp: time.Unix(0, 0).Unix(),
				},
			},
		}}
	r.StoreErrors("test", errors)
	r.StoreErrors("test2", errors)
	services := r.GetServices()
	if !reflect.DeepEqual(services, []string{"test", "test2"}) {
		t.Errorf("Error fetching services,  got %v", services)
	}
}

func TestResolvedErrors(t *testing.T) {
	r := NewORMRepository()
	errors := []ErrorAggregate{
		{
			AggregationKey: "key",
			Severity:       "error",
			CreatedAt:      time.Unix(0, 0).Unix(),
			LatestErrors: []ErrorWithContext{
				{
					Error:     ErrorInstance{},
					Severity:  "error",
					Timestamp: time.Unix(0, 0).Unix(),
				},
			},
		}}
	r.StoreErrors("test", errors)
	r.ResolveError("test", "key")
	if r.countErrors("test") != 0 {
		t.Errorf("Found %d errors, expected 0", r.countErrors("test"))
	}
}

func TestORMRemoveResolved(t *testing.T) {
	r := NewORMRepository()
	errors := []ErrorAggregate{
		{
			AggregationKey: "key",
			Severity:       "error",
			CreatedAt:      time.Unix(0, 0).Unix(),
			LatestErrors: []ErrorWithContext{
				{
					Error:     ErrorInstance{},
					Severity:  "error",
					Timestamp: time.Unix(0, 0).Unix(),
				},
			},
		}}
	r.StoreErrors("test_remove_resolved", errors)
	r.ResolveError("test_remove_resolved", "key")
	r.RemoveResolved("test_remove_resolved", "key")
	if r.countErrors("test_remove_resolved") != 1 {
		t.Errorf("Found %d errors, expected 1", r.countErrors("test_remove_resolved"))
	}
}

func TestORMSearchResolved(t *testing.T) {
	r := NewORMRepository()
	errors := []ErrorAggregate{
		{
			AggregationKey: "key",
			Severity:       "error",
			CreatedAt:      time.Unix(0, 0).Unix(),
			LatestErrors: []ErrorWithContext{
				{
					Error:     ErrorInstance{},
					Severity:  "error",
					Timestamp: time.Unix(0, 0).Unix(),
				},
			},
		}}
	r.StoreErrors("test", errors)
	r.ResolveError("test", "key")

	if !r.SearchResolved("test", "key") {
		t.Errorf("Error should be mark as resolved")
	}
}
