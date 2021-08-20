package repository

import (
	"log"

	"github.com/soundcloud/periskop/config"
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

type Target struct {
	Endpoint string `json:"endpoint"`
}

func NewRepository(repositoryConfig config.Repository) ErrorsRepository {
	if repositoryConfig.Type == "sqlite" {
		log.Printf("Using SQLite %s repository", repositoryConfig.Path)
		db, err := gorm.Open(sqlite.Open(repositoryConfig.Path), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		return NewORMRepository(db)
	} else {
		log.Printf("Using in memory repository")
		return NewMemoryRepository()
	}

}
