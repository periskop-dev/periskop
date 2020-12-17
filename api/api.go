package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/soundcloud/periskop/metrics"
	"github.com/soundcloud/periskop/repository"
)

func NewServicesListHandler(r *repository.ErrorsRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		err := servicesList(w, r)
		if err != nil {
			metrics.ErrorCollector.ReportWithHTTPRequest(err, req)
		}
	})
}

func NewErrorsListHandler(r *repository.ErrorsRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		numberOfOccurrencesPerError := 100

		if service, found := vars["service_name"]; found {
			err := errorsForService(w, r, service, numberOfOccurrencesPerError)
			if err != nil {
				metrics.ErrorCollector.ReportWithHTTPRequest(err, req)
			}
		} else {
			http.NotFound(w, req)
		}
	})
}

func NewErrorResolveHandler(r *repository.ErrorsRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		if service, found := vars["service_name"]; found {
			errKey := vars["error_key"]
			err := (*r).ResolveError(service, errKey)
			if err != nil {
				http.NotFound(w, req)
			}
			w.WriteHeader(http.StatusNoContent)
		} else {
			http.NotFound(w, req)
		}
	})
}

// CORSLocalhostMiddleware allows CORS requests for local development since API and frontend run on different ports
func CORSLocalhostMiddleware(r *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			origin := req.Header.Get("Origin")
			if strings.HasPrefix(origin, "http://localhost:") {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
			}
			next.ServeHTTP(w, req)
		})
	}
}

func errorsForService(w http.ResponseWriter, r *repository.ErrorsRepository,
	service string, numberOfOccurrencesPerError int) error {
	repoErrors, err := (*r).GetErrors(service, numberOfOccurrencesPerError)
	if err == nil {
		err = renderJSON(w, repoErrors)
	} else {
		metrics.ServiceErrors.WithLabelValues("get_errors").Inc()
		http.Error(w, err.Error(), 404)
	}
	return err
}

func servicesList(w http.ResponseWriter, r *repository.ErrorsRepository) error {
	return renderJSON(w, (*r).GetServices())
}

func renderJSON(w http.ResponseWriter, value interface{}) error {
	valueJSON, err := json.Marshal(value)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(valueJSON))
	} else {
		metrics.ServiceErrors.WithLabelValues("render_json").Inc()
		http.Error(w, err.Error(), 500)
	}
	return err
}
