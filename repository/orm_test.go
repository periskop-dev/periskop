package repository

import (
	"reflect"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func countErrors(db *gorm.DB, serviceName string) int64 {
	var count int64
	db.Model(&AggregatedError{}).Where("service_name = ?", serviceName).Count(&count)
	return count
}

func newSQLiteMemory() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(""), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func TestORMReplaceErrors(t *testing.T) {
	db := newSQLiteMemory()
	r := NewORMRepository(db)
	key := "errorKey"
	serviceName := "test_replace"

	// test if creation is successful
	errors := []ErrorAggregate{
		{
			AggregationKey: key,
			Severity:       "error",
			CreatedAt:      time.Unix(0, 0).Unix(),
			TotalCount:     1,
			LatestErrors: []ErrorWithContext{
				{
					Error:     ErrorInstance{},
					Severity:  "error",
					Timestamp: time.Unix(0, 0).Unix(),
				},
			},
		}}
	r.ReplaceErrors(serviceName, errors)
	if countErrors(db, serviceName) != 1 {
		t.Errorf("Found %d errors, expected 1", countErrors(db, serviceName))
	}

	// test if update is successful
	errors = []ErrorAggregate{
		{
			AggregationKey: key,
			Severity:       "error",
			CreatedAt:      time.Unix(0, 0).Unix(),
			TotalCount:     2,
			LatestErrors: []ErrorWithContext{
				{
					Error:     ErrorInstance{},
					Severity:  "error",
					Timestamp: time.Unix(0, 0).Unix(),
				},
			},
		}}
	r.ReplaceErrors(serviceName, errors)
	if countErrors(db, serviceName) != 1 {
		t.Errorf("Found %d errors, expected 1", countErrors(db, serviceName))
	}

	errObj := AggregatedError{}
	db.Model(&AggregatedError{}).
		Where("service_name = ?", serviceName).
		Where("aggregation_key = ?", key).
		First(&errObj)
	if errObj.TotalCount != 2 {
		t.Errorf("Found %d error instances, expected 2", errObj.TotalCount)
	}
}

func TestORMGetErrors(t *testing.T) {
	db := newSQLiteMemory()
	r := NewORMRepository(db)
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

	r.ReplaceErrors("test_get", errors)

	aggregatedErrors, err := r.GetErrors("test_get", 5)
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

func TestORMGetServices(t *testing.T) {
	db := newSQLiteMemory()
	r := NewORMRepository(db)
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
	r.ReplaceErrors("test_services0", errors)
	r.ReplaceErrors("test_services1", errors)
	services := r.GetServices()
	if !reflect.DeepEqual(services, []string{"test_services0", "test_services1"}) {
		t.Errorf("Error fetching services,  got %v", services)
	}
}

func TestORMResolvedErrors(t *testing.T) {
	db := newSQLiteMemory()
	r := NewORMRepository(db)
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
	r.ReplaceErrors("test_resolved", errors)
	r.ResolveError("test_resolved", "key")
	if countErrors(db, "test_resolved") != 0 {
		t.Errorf("Found %d errors, expected 0", countErrors(db, "test_resolved"))
	}
}

func TestORMRemoveResolved(t *testing.T) {
	db := newSQLiteMemory()
	r := NewORMRepository(db)
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
	r.ReplaceErrors("test_remove_resolved", errors)
	r.ResolveError("test_remove_resolved", "key")
	r.RemoveResolved("test_remove_resolved", "key")
	if countErrors(db, "test_remove_resolved") != 1 {
		t.Errorf("Found %d errors, expected 1", countErrors(db, "test_remove_resolved"))
	}
}

func TestORMSearchResolved(t *testing.T) {
	db := newSQLiteMemory()
	r := NewORMRepository(db)
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
	r.ReplaceErrors("test_search", errors)
	r.ReplaceErrors("test_search_other", errors)
	r.ResolveError("test_search", "key")

	if !r.SearchResolved("test_search", "key") {
		t.Errorf("Error should be mark as resolved")
	}

	if r.SearchResolved("test_search_other", "key") {
		t.Errorf("Error shouldn't be mark as resolved")
	}
}
