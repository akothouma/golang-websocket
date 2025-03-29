package tests

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	// "learn.zone01kisumu.ke/git/clomollo/forum/internal/Handlers"
)

type MockDependencies struct {
	ErrorLog *log.Logger
}

func (m *MockDependencies) ServerError(w http.ResponseWriter, err error) {
	trace := err.Error()
	m.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (m *MockDependencies) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (m *MockDependencies) NotFound(w http.ResponseWriter) {
	m.ClientError(w, http.StatusNotFound)
}

func TestServerError(t *testing.T) {
	logBuffer := new(bytes.Buffer)
	mockDeps := &MockDependencies{ErrorLog: log.New(logBuffer, "", log.LstdFlags)}

	// req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()
	mockDeps.ServerError(rr, errors.New("Test Server Error"))

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}

	if logBuffer.Len() == 0 {
		t.Errorf("Expected log output but got none")
	}
}

func TestClientError(t *testing.T) {
	mockDeps := &MockDependencies{}
	// req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	mockDeps.ClientError(rr, http.StatusBadRequest)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestNotFound(t *testing.T) {
	mockDeps := &MockDependencies{}
	// req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	mockDeps.NotFound(rr)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, rr.Code)
	}
}
