package tests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/Handlers"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

func TestAddCommentHandler(t *testing.T) {
	originalAddComment := models.AddComment
	defer func() { models.AddComment = originalAddComment }()

	tests := []struct {
		name           string
		method         string
		setupRequest   func() *http.Request
		mockAddComment func(string, string, string) (int64, error)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Non-POST method",
			method: http.MethodGet,
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/add-comment", nil)
				req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session"})
				ctx := context.WithValue(req.Context(), "session_id", "test-session")
				ctx = context.WithValue(ctx, "user_uuid", "test-user")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "Invalid request method",
		},
		{
			name:   "Missing session cookie",
			method: http.MethodPost,
			setupRequest: func() *http.Request {
				form := strings.NewReader("post_id=123&content=test")
				req := httptest.NewRequest(http.MethodPost, "/add-comment", form)
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				ctx := context.WithValue(req.Context(), "session_id", "test-session")
				ctx = context.WithValue(ctx, "user_uuid", "test-user")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name:   "Session cookie mismatch",
			method: http.MethodPost,
			setupRequest: func() *http.Request {
				form := strings.NewReader("post_id=123&content=test")
				req := httptest.NewRequest(http.MethodPost, "/add-comment", form)
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				req.AddCookie(&http.Cookie{Name: "session_id", Value: "cookie-session"})
				ctx := context.WithValue(req.Context(), "session_id", "context-session")
				ctx = context.WithValue(ctx, "user_uuid", "test-user")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name:   "Empty content",
			method: http.MethodPost,
			setupRequest: func() *http.Request {
				form := strings.NewReader("post_id=123&content=")
				req := httptest.NewRequest(http.MethodPost, "/add-comment", form)
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session"})
				ctx := context.WithValue(req.Context(), "session_id", "test-session")
				ctx = context.WithValue(ctx, "user_uuid", "test-user")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Missing required fields",
		},
		{
			name:   "Success",
			method: http.MethodPost,
			setupRequest: func() *http.Request {
				form := strings.NewReader("post_id=123&content=test+comment")
				req := httptest.NewRequest(http.MethodPost, "/add-comment", form)
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session"})
				ctx := context.WithValue(req.Context(), "session_id", "test-session")
				ctx = context.WithValue(ctx, "user_uuid", "test-user")
				return req.WithContext(ctx)
			},
			mockAddComment: func(postID, userID, content string) (int64, error) {
				if postID != "123" || userID != "test-user" || content != "test comment" {
					return 0, fmt.Errorf("unexpected arguments")
				}
				return 456, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"comment_id":456,"message":"Comment added successfully"}`,
		},
		{
			name:   "AddComment error",
			method: http.MethodPost,
			setupRequest: func() *http.Request {
				form := strings.NewReader("post_id=123&content=test+comment")
				req := httptest.NewRequest(http.MethodPost, "/add-comment", form)
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session"})
				ctx := context.WithValue(req.Context(), "session_id", "test-session")
				ctx = context.WithValue(ctx, "user_uuid", "test-user")
				return req.WithContext(ctx)
			},
			mockAddComment: func(postID, userID, content string) (int64, error) {
				return 0, fmt.Errorf("database error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to add comment: database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockAddComment != nil {
				models.AddComment = tt.mockAddComment
			} else {
				models.AddComment = originalAddComment
			}

			req := tt.setupRequest()
			rr := httptest.NewRecorder()

			handlers.AddCommentHandler(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			body := strings.TrimSpace(rr.Body.String())
			if tt.expectedBody != "" && body != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, body)
			}

			if tt.expectedStatus == http.StatusOK && tt.mockAddComment != nil {
				if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
					t.Errorf("expected Content-Type application/json, got %s", contentType)
				}
			}
		})
	}
}
