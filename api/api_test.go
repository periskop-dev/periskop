package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/soundcloud/periskop/repository"
)

func TestServicesWithEmptyRepoReturnsSuccess(t *testing.T) {
	r := repository.NewInMemory()

	rr := httptest.NewRecorder()
	serveMockServiceList(rr, r)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestServicesWithEmptyRepoReturnsEmptyArray(t *testing.T) {
	r := repository.NewInMemory()

	rr := httptest.NewRecorder()
	serveMockServiceList(rr, r)

	expected := "[]\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestServicesWithNonEmptyRepoReturnsServiceNames(t *testing.T) {
	r := repository.NewInMemory()
	r.StoreErrors("api-test", []repository.ErrorAggregate{})

	rr := httptest.NewRecorder()
	serveMockServiceList(rr, r)

	expected := "[\"api-test\"]\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func serveMockServiceList(rr *httptest.ResponseRecorder, r repository.ErrorsRepository) {
	handler := NewServicesListHandler(&r)
	router := mux.NewRouter()
	router.Handle("/services/", handler).Methods(http.MethodGet)
	req, _ := http.NewRequest("GET", "/services/", nil)
	router.ServeHTTP(rr, req)
}

func TestErrorsForUnknownServiceReturnsNotFound(t *testing.T) {
	r := repository.NewInMemory()
	rr := httptest.NewRecorder()
	serveMockErrorList(rr, r, "api-test")

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestErrorsForKnownServiceReturnsSuccess(t *testing.T) {
	r := repository.NewInMemory()
	r.StoreErrors("api-test", []repository.ErrorAggregate{})

	rr := httptest.NewRecorder()
	serveMockErrorList(rr, r, "api-test")

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestErrorsForKnownServiceReturnsErrors(t *testing.T) {
	r := repository.NewInMemory()
	r.StoreErrors("api-test", []repository.ErrorAggregate{
		{
			AggregationKey: "key",
			Severity:       "error",
			CreatedAt:      time.Unix(0, 0).Unix(),
			LatestErrors: []repository.ErrorWithContext{
				{
					Error:     repository.ErrorInstance{},
					Severity:  "error",
					Timestamp: time.Unix(0, 0).Unix(),
				},
			},
		}})

	rr := httptest.NewRecorder()
	serveMockErrorList(rr, r, "api-test")

	// nolint
	expected := `[{"aggregation_key":"key","total_count":0,"severity":"error","latest_errors":[{"error":{"class":"","message":"","stacktrace":null,"cause":null},"uuid":"","timestamp":0,"severity":"error","http_context":null}],"created_at":0}]` + "\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func serveMockErrorList(rr *httptest.ResponseRecorder, r repository.ErrorsRepository, serviceName string) {
	handler := NewErrorsListHandler(&r)
	router := mux.NewRouter()
	router.Handle("/services/{service_name}/errors/", handler).Methods(http.MethodGet)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/services/%s/errors/", serviceName), nil)
	router.ServeHTTP(rr, req)
}

func TestResolveErrorForUnknownServiceReturnsNotFound(t *testing.T) {
	r := repository.NewInMemory()
	rr := httptest.NewRecorder()
	serveMockErrorResolve(rr, r, "api-test", "test")

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestResolveErrorsReturnsSuccess(t *testing.T) {
	r := repository.NewInMemory()
	r.StoreErrors("api-test", []repository.ErrorAggregate{})

	rr := httptest.NewRecorder()
	serveMockErrorResolve(rr, r, "api-test", "test")

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func serveMockErrorResolve(rr *httptest.ResponseRecorder, r repository.ErrorsRepository,
	serviceName string, errKey string) {
	handler := NewErrorResolveHandler(&r)
	router := mux.NewRouter()
	router.Handle("/services/{service_name}/errors/{error_key}/", handler).Methods(http.MethodDelete)
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/services/%s/errors/%s/", serviceName, errKey), nil)
	router.ServeHTTP(rr, req)
}
