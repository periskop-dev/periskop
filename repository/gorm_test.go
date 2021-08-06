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
