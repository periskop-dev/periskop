package repository

import (
	"log"

	"github.com/soundcloud/periskop/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ErrorAggregate struct {
	AggregationKey string             `json:"aggregation_key"`
	TotalCount     int                `json:"total_count"`
	Severity       string             `json:"severity"`
	LatestErrors   []ErrorWithContext `json:"latest_errors"`
	CreatedAt      int64              `json:"created_at"`
}

type ErrorWithContext struct {
	Error       ErrorInstance `json:"error"`
	UUID        string        `json:"uuid"`
	Timestamp   int64         `json:"timestamp"`
	Severity    string        `json:"severity"`
	HTTPContext *HTTPContext  `json:"http_context"`
}

type ErrorInstance struct {
	Class      string         `json:"class"`
	Message    string         `json:"message"`
	Stacktrace []string       `json:"stacktrace"`
	Cause      *ErrorInstance `json:"cause"`
}

type HTTPContext struct {
	RequestMethod  string            `json:"request_method"`
	RequestURL     string            `json:"request_url"`
	RequestHeaders map[string]string `json:"request_headers"`
	RequestBody    string            `json:"request_body"`
}

type Target struct {
	Endpoint string `json:"endpoint"`
}

type ErrorsRepository interface {
	GetErrors(serviceName string, numberOfErrors int) ([]ErrorAggregate, error)
	StoreErrors(serviceName string, errors []ErrorAggregate)
	GetServices() []string
	ResolveError(serviceName string, key string) error
	SearchResolved(serviceName string, key string) bool
	RemoveResolved(serviceName string, key string)
	StoreTargets(serviceName string, targets []Target)
	GetTargets() map[string][]Target
}

// NewRepository it's a factory function for ErrorRepository interfaces.
// It creates a repository based on the configured repository.
func NewRepository(repositoryConfig config.Repository) ErrorsRepository {
	if repositoryConfig.Type == "sqlite" {
		log.Printf("Using SQLite %s repository", repositoryConfig.Path)
		db, err := gorm.Open(sqlite.Open(repositoryConfig.Path), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		return NewORMRepository(db)
	} else if repositoryConfig.Type == "mysql" {
		log.Printf("Using MySQL repository")
		db, err := gorm.Open(mysql.Open(repositoryConfig.Dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		return NewORMRepository(db)
	} else if repositoryConfig.Type == "postgres" {
		log.Printf("Using PostgresSQL repository")
		db, err := gorm.Open(postgres.Open(repositoryConfig.Dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		return NewORMRepository(db)
	} else {
		log.Printf("Using in memory repository")
		return NewMemoryRepository()
	}
}
