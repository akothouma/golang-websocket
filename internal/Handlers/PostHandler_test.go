package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostsByFilters(t *testing.T) {
    categoriesForFilter:=[]string{"education"}
	validData := []byte({"categories":$categoriesForFilter})
	req := httptest.NewRequest(http.MethodPost, "/filtered_posts", bytes.NewBuffer(validData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	PostsByFilters(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, w.Code)
	}
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type %v, got %v", "application/json", contentType)
	}
	var response []string
	json.Unmarshal(w.Body.Bytes(), &response)
    
	if len(response) == 0 {
        t.Error("Expected at least one post, got empty slice")
    }

	req = httptest.NewRequest(http.MethodGet, "/filtered_posts", nil)
	w = httptest.NewRecorder()
	PostsByFilters(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %v, got %v", http.StatusMethodNotAllowed, w.Code)
	}

	invalidJSON := []byte(`{"categories": categoriesForFilter`)
	req = httptest.NewRequest(http.MethodPost, "/filtered_posts", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	PostsByFilters(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %v, got %v", http.StatusBadRequest, w.Code)
	}

	missingData := []byte(`{"categories":}`)
	req = httptest.NewRequest(http.MethodPost, "/filtered_posts", bytes.NewBuffer(missingData))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	PostsByFilters(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %v, got %v", http.StatusBadRequest, w.Code)
	}
}
