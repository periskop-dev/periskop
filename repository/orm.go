package repository

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/soundcloud/periskop/metrics"
	"gorm.io/gorm"
)

type ormRepository struct {
	DB *gorm.DB
	targetsRepository
}

func (e *ErrorAggregate) Scan(src interface{}) error {
	return json.Unmarshal([]byte(src.(string)), &e)
}

func (e ErrorAggregate) Value() (driver.Value, error) {
	val, err := json.Marshal(e)
	return string(val), err
}

type AggregatedError struct {
	gorm.Model
	ServiceName    string `gorm:"index"`
	AggregationKey string `gorm:"index"`
	Errors         ErrorAggregate
	TotalCount     int
}

func NewORMRepository(db *gorm.DB) ErrorsRepository {
	err := db.AutoMigrate(&AggregatedError{})
	if err != nil {
		panic("failed to create database migration")
	}
	return &ormRepository{DB: db}
}

func (r *ormRepository) GetErrors(serviceName string, numberOfErrors int) ([]ErrorAggregate, error) {
	aggregatedErrors := []AggregatedError{}
	r.DB.
		Where(&AggregatedError{ServiceName: serviceName}).
		Find(&aggregatedErrors)

	errors := []ErrorAggregate{}
	for _, aggregatedError := range aggregatedErrors {
		errorObj := aggregatedError.Errors
		maxErrors := len(errorObj.LatestErrors)
		if numberOfErrors < maxErrors {
			maxErrors = numberOfErrors
		}
		errorObj.LatestErrors = errorObj.LatestErrors[0:maxErrors]
		errors = append(errors, errorObj)
	}
	if len(errors) > 0 {
		return errors, nil
	}
	metrics.ServiceErrors.WithLabelValues("service_not_found").Inc()
	return nil, fmt.Errorf("service %s not found", serviceName)
}

// ReplaceErrors deletes previous stored errors for a service name and stores the new list of errors in json format
func (r *ormRepository) ReplaceErrors(serviceName string, errors []ErrorAggregate) {
	for _, errorAggregate := range errors {
		errObj := AggregatedError{}
		key := errorAggregate.AggregationKey
		result := r.DB.Model(&AggregatedError{}).
			Where("service_name = ?", serviceName).
			Where("aggregation_key = ?", key).
			First(&errObj)
		if result.RowsAffected == 0 {
			r.DB.Create(&AggregatedError{
				ServiceName:    serviceName,
				Errors:         errorAggregate,
				AggregationKey: key,
				TotalCount:     errorAggregate.TotalCount,
			})
		} else if errorAggregate.TotalCount > errObj.TotalCount { // only update if there are more errors than before
			r.DB.Model(&AggregatedError{}).
				Where("service_name = ?", serviceName).
				Where("aggregation_key = ?", key).
				Update("total_count", errorAggregate.TotalCount).
				Update("errors", errorAggregate)
		}
	}
}

// GetServices fetches the list of unique services
func (r *ormRepository) GetServices() []string {
	aggregatedErrors := []AggregatedError{}
	keys := make([]string, 0)
	r.DB.
		Distinct("service_name").
		Find(&aggregatedErrors)
	for _, aggregatedError := range aggregatedErrors {
		keys = append(keys, aggregatedError.ServiceName)
	}
	return keys
}

// ResolveError marks the given error as soft-deleted (updating deleted_at)
func (r *ormRepository) ResolveError(serviceName string, key string) error {
	r.DB.
		Where("service_name = ?", serviceName).
		Where("aggregation_key = ?", key).
		Delete(&AggregatedError{})
	return nil
}

// RemoveResolved removes a soft-deletion of the given error
func (r *ormRepository) RemoveResolved(serviceName string, key string) {
	r.DB.Model(&AggregatedError{}).
		Where("service_name = ?", serviceName).
		Where("aggregation_key = ?", key).
		Unscoped().
		Update("deleted_at", nil)
}

// SearchResolved returns true if the given error was marked previously as resolved
func (r *ormRepository) SearchResolved(serviceName string, key string) bool {
	var count int64
	r.DB.
		Model(&AggregatedError{}).
		Where("service_name = ?", serviceName).
		Where("aggregation_key = ?", key).
		Where("deleted_at is NOT NULL").
		Unscoped().
		Count(&count)
	return count >= 1
}
