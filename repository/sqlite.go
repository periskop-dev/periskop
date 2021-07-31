package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/soundcloud/periskop/metrics"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository() ErrorsRepository {
	db, _ := sql.Open("sqlite3", "./periskop.db")
	createTables(db)
	return &SQLiteRepository{db}
}

func createTables(db *sql.DB) {
	createErrorsTable := `CREATE TABLE error (
		"service_name" TEXT,
		"aggregation_key" TEXT,
		"error_json" TEXT		
	  );`

	log.Println("Create errors table...")
	statement, err := db.Prepare(createErrorsTable)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
}

func (r *SQLiteRepository) GetErrors(serviceName string, numberOfErrors int) ([]ErrorAggregate, error) {
	row, err := r.db.Query("SELECT * FROM error WHERE service_name='%s' LIMIT='%s'", serviceName, numberOfErrors)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	if value, ok := r.AggregatedError.Load(serviceName); ok {
		value, _ := value.([]ErrorAggregate)
		result := make([]ErrorAggregate, 0, len(value))
		for _, errorObj := range value {
			topCap := len(errorObj.LatestErrors)
			if numberOfErrors < topCap {
				topCap = numberOfErrors
			}
			errorObj.LatestErrors = errorObj.LatestErrors[0:topCap]
			result = append(result, errorObj)
		}

		return result, nil
	}
	metrics.ServiceErrors.WithLabelValues("service_not_found").Inc()
	return nil, fmt.Errorf("service %s not found", serviceName)
}
