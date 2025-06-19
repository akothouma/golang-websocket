package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/Handlers"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

// Mock function for GetAllCommentsForPost
func mockGetAllCommentsForPost(postID string) ([]models.Comment, error) {
	if postID == "valid_id" {
		return []models.Comment{
			{ID: 1, PostID: "valid_id", Content: "Test Comment", UserUuiD: "user1"},
		}, nil
	}
	return nil, errors.New("database error")
}

func TestGetAllCommentsForPostHandler(t *testing.T) {
	// Backup and restore the original function after test
	originalFunc := models.GetAllCommentsForPost
	models.GetAllCommentsForPost = mockGetAllCommentsForPost
	defer func() { models.GetAllCommentsForPost = originalFunc }()

	// Test cases
	tests := []struct {
		name           string
		postID         string
		expectedStatus int
		expectError    bool
	}{
		{"Valid Request", "valid_id", http.StatusOK, false},
		{"Missing Post ID", "", http.StatusBadRequest, true},
		{"Database Error", "invalid_id", http.StatusInternalServerError, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/comments?post_id="+test.postID, nil)
			rr := httptest.NewRecorder()
			handlers.GetAllCommentsForPostHandler(rr, req)

			if rr.Code != test.expectedStatus {
				t.Errorf("Expected status %d, got %d", test.expectedStatus, rr.Code)
			}

			if !test.expectError {
				var comments []models.Comment
				if err := json.Unmarshal(rr.Body.Bytes(), &comments); err != nil {
					t.Errorf("Failed to parse JSON response: %v", err)
				}
			}
		})
	}
}
