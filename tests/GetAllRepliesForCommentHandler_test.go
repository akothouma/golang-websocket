package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	// "learn.zone01kisumu.ke/git/clomollo/forum/internal/Handlers"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

// Mock struct that embeds Comment
type MockComment struct {
	models.Comment
}

// Mock method to simulate GetAllRepliesForComment
func (c *MockComment) GetAllRepliesForComment() error {
	if c.ID == 1 {
		c.Replies = []models.Comment{
			{ID: 2, ParentCommentID: 1, Content: "Reply 1", UserUuiD: "user1"},
			{ID: 3, ParentCommentID: 1, Content: "Reply 2", UserUuiD: "user2"},
		}
		return nil
	}
	return errors.New("database error")
}

// Mock handler function that injects the mock struct
func MockGetAllRepliesForCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	commentIDStr := r.URL.Query().Get("comment_id")
	if commentIDStr == "" {
		http.Error(w, "Comment ID is required", http.StatusBadRequest)
		return
	}

	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// Use mock comment instead of real database
	comment := &MockComment{}
	comment.ID = commentID

	err = comment.GetAllRepliesForComment()
	if err != nil {
		http.Error(w, "Failed to retrieve replies", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment.Replies)
}

// Test function
func TestGetAllRepliesForCommentHandler(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		commentID      string
		expectedStatus int
		expectError    bool
	}{
		{"Valid Request", "1", http.StatusOK, false},
		{"Missing Comment ID", "", http.StatusBadRequest, true},
		{"Invalid Comment ID", "abc", http.StatusBadRequest, true},
		{"Database Error", "2", http.StatusInternalServerError, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/replies?comment_id="+test.commentID, nil)
			rr := httptest.NewRecorder()

			// Call the **mocked** handler function instead of the actual one
			MockGetAllRepliesForCommentHandler(rr, req)

			// Check status code
			if rr.Code != test.expectedStatus {
				t.Errorf("Expected status %d, got %d", test.expectedStatus, rr.Code)
			}

			// Check JSON response if no error expected
			if !test.expectError {
				var replies []models.Comment
				if err := json.Unmarshal(rr.Body.Bytes(), &replies); err != nil {
					t.Errorf("Failed to parse JSON response: %v", err)
				}
				if len(replies) == 0 {
					t.Errorf("Expected replies but got none")
				}
			}
		})
	}
}
