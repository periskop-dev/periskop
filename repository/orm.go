package repository

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/soundcloud/periskop/metrics"
	"gorm.io/gorm"
)

type ormRepository struct {
	DB *gorm.DB
	// map service name -> list of scraped targets
	Targets sync.Map
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
	ServiceName    string
	AggregationKey string
	Errors         ErrorAggregate
}

func NewORMRepository(db *gorm.DB) ErrorsRepository {
	db.AutoMigrate(&AggregatedError{})
	return &ormRepository{DB: db, Targets: sync.Map{}}
}

func (r *ormRepository) GetErrors(serviceName string, numberOfErrors int) ([]ErrorAggregate, error) {
	aggregatedErrors := []AggregatedError{}
	r.DB.
		Where(&AggregatedError{ServiceName: serviceName}).
		Find(&aggregatedErrors)

	result := []ErrorAggregate{}
	for _, aggregatedError := range aggregatedErrors {
		errorObj := aggregatedError.Errors
		topCap := len(errorObj.LatestErrors)
		if numberOfErrors < topCap {
			topCap = numberOfErrors
		}
		errorObj.LatestErrors = errorObj.LatestErrors[0:topCap]
		result = append(result, errorObj)
	}
	if len(result) > 0 {
		return result, nil
	} else {
		metrics.ServiceErrors.WithLabelValues("service_not_found").Inc()
		return nil, fmt.Errorf("service %s not found", serviceName)
	}
}

// StoreErrors deletes previous stored errors for a service name and store the new list of errors in json format
func (r *ormRepository) StoreErrors(serviceName string, errors []ErrorAggregate) {
	// Delete previous records
	r.DB.
		Where("service_name = ?", serviceName).
		Unscoped().
		Delete(&AggregatedError{})

	for _, errorAggregate := range errors {
		r.DB.Create(&AggregatedError{
			ServiceName:    serviceName,
			Errors:         errorAggregate,
			AggregationKey: errorAggregate.AggregationKey,
		})
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
	return count == 1
}

func (r *ormRepository) StoreTargets(serviceName string, targets []Target) {
	r.Targets.Store(serviceName, targets)
}

func (r *ormRepository) GetTargets() map[string][]Target {
	targets := make(map[string][]Target)
	r.Targets.Range(func(key, value interface{}) bool {
		targets[key.(string)] = value.([]Target)
		return true
	})
	return targets
}
