package tests

import (
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/Handlers"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestErrorHandler(t *testing.T) {
	// Ensure the directory exists
	err := os.MkdirAll("./ui/html", os.ModePerm)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create a mock error template
	templatePath := "./ui/html/error.html"
	templateContent := "<html><body>Error: {{.Code}}</body></html>"
	err = os.WriteFile(templatePath, []byte(templateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create template file: %v", err)
	}
	defer os.Remove(templatePath) // Clean up after test

	// Create test cases
	tests := []struct {
		code     int
		expected string
	}{
		{http.StatusNotFound, "Error: 404"},
		{http.StatusInternalServerError, "Error: 500"},
	}

	for _, test := range tests {
		rr := httptest.NewRecorder()
		handlers.ErrorHandler(rr, test.code)

		// Check status code
		if rr.Code != test.code {
			t.Errorf("Expected status %d, got %d", test.code, rr.Code)
		}

		// Check response body
		if !strings.Contains(rr.Body.String(), test.expected) {
			t.Errorf("Expected body to contain '%s', got '%s'", test.expected, rr.Body.String())
		}
	}
}
