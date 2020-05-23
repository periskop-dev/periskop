package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/soundcloud/periskop/metrics"
	"github.com/soundcloud/periskop/repository"
)

func NewHandler(r *repository.ErrorsRepository, serverURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Allow CORS requests for local development since API and frontend run on different ports
		origin := req.Header.Get("Origin")
		if strings.HasPrefix(origin, fmt.Sprintf("http://%s:", serverURL)) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		path := strings.TrimPrefix(req.URL.Path, "/services/")
		numberOfOccurrencesPerError := 10

		if len(path) == 0 {
			err := servicesList(w, r)
			if err != nil {
				metrics.ErrorCollector.ReportWithHTTPRequest(err, req)
			}
		} else if service, err := extractServiceName(path); err == nil {
			err = errorsForService(w, r, service, numberOfOccurrencesPerError)
			if err != nil {
				metrics.ErrorCollector.ReportWithHTTPRequest(err, req)
			}
		} else {
			metrics.ErrorCollector.ReportWithHTTPRequest(err, req)
			http.NotFound(w, req)
		}
	})
}

func extractServiceName(url string) (string, error) {
	if strings.HasSuffix(url, "/errors/") {
		return strings.TrimSuffix(url, "/errors/"), nil
	}
	return "", errors.New("invalid path")
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
