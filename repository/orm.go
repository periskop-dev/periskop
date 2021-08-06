package repository

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/soundcloud/periskop/metrics"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ormRepository struct {
	db *gorm.DB
}

type ErrorsRepository2 interface {
	StoreErrors(serviceName string, errors []ErrorAggregate)
	GetErrors(serviceName string, numberOfErrors int) ([]ErrorAggregate, error)
	countErrors(serviceName string) int64
	ResolveError(serviceName string, key string) error
	RemoveResolved(serviceName string, key string)
	SearchResolved(serviceName string, key string) bool
	GetServices() []string
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

func NewORMRepository() ErrorsRepository2 {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&AggregatedError{})
	return &ormRepository{db}
}

func (r *ormRepository) countErrors(serviceName string) int64 {
	var count int64
	r.db.Model(&AggregatedError{}).Where("service_name = ?", serviceName).Count(&count)
	return count
}

func (r *ormRepository) GetErrors(serviceName string, numberOfErrors int) ([]ErrorAggregate, error) {
	aggregatedErrors := []AggregatedError{}
	r.db.Where(&AggregatedError{ServiceName: serviceName}).Find(&aggregatedErrors)

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

func (r *ormRepository) StoreErrors(serviceName string, errors []ErrorAggregate) {
	// Delete previous records
	r.db.Where("service_name = ?", serviceName).Unscoped().Delete(&AggregatedError{})
	//r.db.Exec("DELETE FROM aggregated_errors where service_name = ?", serviceName)

	for _, errorAggregate := range errors {
		r.db.Create(&AggregatedError{
			ServiceName:    serviceName,
			Errors:         errorAggregate,
			AggregationKey: errorAggregate.AggregationKey,
		})
	}
}

func (r *ormRepository) GetServices() []string {
	aggregatedErrors := []AggregatedError{}
	keys := make([]string, 0)
	r.db.Distinct("service_name").Find(&aggregatedErrors)
	for _, aggregatedError := range aggregatedErrors {
		keys = append(keys, aggregatedError.ServiceName)
	}
	return keys
}

func (r *ormRepository) ResolveError(serviceName string, key string) error {
	r.db.Where("service_name = ?", serviceName).Where("aggregation_key = ?", key).Delete(&AggregatedError{})
	return nil
}

func (r *ormRepository) RemoveResolved(serviceName string, key string) {
	r.db.Model(&AggregatedError{}).Where("service_name = ?", serviceName).Where("aggregation_key = ?", key).Update("deleted_at", "null")
}

func (r *ormRepository) SearchResolved(serviceName string, key string) bool {
	var count int64
	r.db.Model(&AggregatedError{ServiceName: serviceName, AggregationKey: key}).Unscoped().Count(&count)
	return count == 1
}
