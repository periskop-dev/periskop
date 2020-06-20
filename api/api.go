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
		allowCORS(w, req)
		err := servicesList(w, r)
		if err != nil {
			metrics.ErrorCollector.ReportWithHTTPRequest(err, req)
		}
	})
}

func NewErrorsListHandler(r *repository.ErrorsRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		allowCORS(w, req)
		vars := mux.Vars(req)
		numberOfOccurrencesPerError := 10

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

func allowCORS(w http.ResponseWriter, req *http.Request) {
	// Allow CORS requests for local development since API and frontend run on different ports
	origin := req.Header.Get("Origin")
	if strings.HasPrefix(origin, "http://localhost:") {
		w.Header().Set("Access-Control-Allow-Origin", origin)
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
