package repository

import (
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
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
}
