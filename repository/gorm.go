package repository

import (
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
}

type AggregatedError struct {
	gorm.Model
	ServiceName string
	Errors      []ErrorAggregate
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

func (r *ormRepository) GetErrors(serviceName string, numberOfErrors int) ([]ErrorAggregate, error) {
	aggr := AggregatedError{}
	r.db.Where("service_name = ?", serviceName).Find(&aggr)
	if len(aggr.Errors) > 0 {
		topCap := len(aggr.Errors)
		if numberOfErrors < topCap {
			topCap = numberOfErrors
		}
		return aggr.Errors[0:topCap], nil
	} else {
		metrics.ServiceErrors.WithLabelValues("service_not_found").Inc()
		return nil, fmt.Errorf("service %s not found", serviceName)
	}
}
