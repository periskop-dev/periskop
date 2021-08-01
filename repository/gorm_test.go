package repository

import (
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
	if r.countErrors("test") == 0 {
		t.Errorf("Found 0 errors, expected 1")
	}

	aggr, err := r.GetErrors("test", 5)
	if err != nil {
		t.Errorf("Fail to fetch errors: %s", err)
	}

	if len(aggr) < 0 {
		t.Errorf("Found 0 errors, expected 1")
	}
}
