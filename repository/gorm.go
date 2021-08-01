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
}

type ErrorsArray []ErrorAggregate

func (e *ErrorsArray) Scan(src interface{}) error {
	return json.Unmarshal([]byte(src.(string)), &e)
}

func (e ErrorsArray) Value() (driver.Value, error) {
	val, err := json.Marshal(e)
	return string(val), err
}

type AggregatedError struct {
	gorm.Model
	ServiceName string
	Errors      ErrorsArray
}

func NewORMRepository() ErrorsRepository2 {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&AggregatedError{})
	return &ormRepository{db}
}

func (r *ormRepository) StoreErrors(serviceName string, errors []ErrorAggregate) {
	r.db.Create(&AggregatedError{ServiceName: serviceName, Errors: errors})
}

func (r *ormRepository) countErrors(serviceName string) int64 {
	res := r.db.Find(&AggregatedError{})
	return res.RowsAffected
}

func (r *ormRepository) GetErrors(serviceName string, numberOfErrors int) ([]ErrorAggregate, error) {
	aggregatedErr := AggregatedError{}
	r.db.Where(&AggregatedError{ServiceName: serviceName}).First(&aggregatedErr)
	result := []ErrorAggregate{}
	for _, errorObj := range aggregatedErr.Errors {
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
