package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/soundcloud/periskop/repository"
)

func NewHandler(r repository.ErrorsRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Allow CORS requests for local development since API and frontend run on different ports
		origin := req.Header.Get("Origin")
		if strings.HasPrefix(origin, "http://localhost:") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		path := strings.TrimPrefix(req.URL.Path, "/services/")
		numberOfOccurrencesPerError := 10

		if len(path) == 0 {
			servicesList(w, r)
		} else if service, err := extractServiceName(path); err == nil {
			errorsForService(w, r, service, numberOfOccurrencesPerError)
		} else {
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

func errorsForService(w http.ResponseWriter, r repository.ErrorsRepository, service string, numberOfOccurrencesPerError int) {
	if repoErrors, err := r.GetErrors(service, numberOfOccurrencesPerError); err == nil {
		renderJSON(w, repoErrors)
	} else {
		http.Error(w, err.Error(), 404)
	}
}

func servicesList(w http.ResponseWriter, r repository.ErrorsRepository) {
	renderJSON(w, r.GetServices())
}

func renderJSON(w http.ResponseWriter, value interface{}) {
	if valueJSON, err := json.Marshal(value); err == nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, string(valueJSON))
	} else {
		http.Error(w, err.Error(), 500)
	}
}
