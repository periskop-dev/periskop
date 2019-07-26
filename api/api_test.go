package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/soundcloud/periskop/repository"
)

func TestRandomPathReturnNotFound(t *testing.T) {
	r := repository.NewInMemory()
	handler := NewHandler(r)

	req, _ := http.NewRequest("GET", "/whatever", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestServicesWithEmptyRepoReturnsSuccess(t *testing.T) {
	r := repository.NewInMemory()
	handler := NewHandler(r)

	req, _ := http.NewRequest("GET", "/services/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestServicesWithEmptyRepoReturnsEmptyArray(t *testing.T) {
	r := repository.NewInMemory()
	handler := NewHandler(r)

	req, _ := http.NewRequest("GET", "/services/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	expected := "[]\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestServicesWithNonEmptyRepoReturnsServiceNames(t *testing.T) {
	r := repository.NewInMemory()
	r.StoreErrors("api-test", []repository.ErrorAggregate{})

	handler := NewHandler(r)

	req, _ := http.NewRequest("GET", "/services/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	expected := "[\"api-test\"]\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestErrorsForUnknownServiceReturnsNotFound(t *testing.T) {
	r := repository.NewInMemory()
	handler := NewHandler(r)

	req, _ := http.NewRequest("GET", "/services/api-test/erros/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func TestErrorsForKnownServiceReturnsSuccess(t *testing.T) {
	r := repository.NewInMemory()
	r.StoreErrors("api-test", []repository.ErrorAggregate{})

	handler := NewHandler(r)

	req, _ := http.NewRequest("GET", "/services/api-test/errors/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

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
			LatestErrors: []repository.ErrorWithContext{
				{
					Error:     repository.ErrorInstance{},
					Timestamp: time.Unix(0, 0).Unix(),
				},
			},
		}})

	handler := NewHandler(r)

	req, _ := http.NewRequest("GET", "/services/api-test/errors/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// nolint
	expected := `[{"aggregation_key":"key","total_count":0,"latest_errors":[{"error":{"class":"","message":"","stacktrace":null,"cause":null},"uuid":"","timestamp":0,"http_context":{"request_method":"","request_url":"","request_headers":null}}]}]` + "\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
