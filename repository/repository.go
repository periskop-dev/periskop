package repository

import (
	"log"
	"sync"

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

type TargetsRepository interface {
	StoreTargets(serviceName string, targets []Target)
	GetTargets() map[string][]Target
}

type ErrorsRepository interface {
	GetErrors(serviceName string, numberOfErrors int) ([]ErrorAggregate, error)
	ReplaceErrors(serviceName string, errors []ErrorAggregate)
	GetServices() []string
	ResolveError(serviceName string, key string) error
	SearchResolved(serviceName string, key string) bool
	RemoveResolved(serviceName string, key string)
	TargetsRepository
}

type targetsRepository struct {
	// map service name -> list of scraped targets
	Targets sync.Map
}

// StoreTargets stores a list of scrapped targets (hosts) for a service
func (r *targetsRepository) StoreTargets(serviceName string, targets []Target) {
	r.Targets.Store(serviceName, targets)
}

// GetTargets gets a list of scrapped targets (hosts) for a service
func (r *targetsRepository) GetTargets() map[string][]Target {
	targets := make(map[string][]Target)
	r.Targets.Range(func(key, value interface{}) bool {
		targets[key.(string)] = value.([]Target)
		return true
	})
	return targets
}

// NewRepository is a factory function for ErrorRepository interfaces.
// It creates a repository based on the configured repository.
func NewRepository(repositoryConfig config.Repository) ErrorsRepository {
	switch repositoryConfig.Type {
	case "sqlite":
		log.Printf("Using SQLite %s repository", repositoryConfig.Path)
		db, err := gorm.Open(sqlite.Open(repositoryConfig.Path), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		return NewORMRepository(db)
	case "mysql":
		log.Printf("Using MySQL repository")
		db, err := gorm.Open(mysql.Open(repositoryConfig.Dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		return NewORMRepository(db)
	case "postgres":
		log.Printf("Using PostgresSQL repository")
		db, err := gorm.Open(postgres.Open(repositoryConfig.Dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		return NewORMRepository(db)
	case "memory":
		log.Printf("Using in memory repository")
		return NewMemoryRepository()
	default:
		log.Printf("Using in memory repository")
		return NewMemoryRepository()
	}
}
