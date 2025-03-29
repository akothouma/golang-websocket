package tests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/Handlers"
)

func TestAddReplyHandler(t *testing.T) {
	originalAddReply := models.AddReply
	defer func() { models.AddReply = originalAddReply }()

	tests := []struct {
		name           string
		method         string
		setupRequest   func() *http.Request
		mockAddReply   func(string, string, string) (int64, error)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Non-POST method",
			method: http.MethodGet,
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/add-reply", nil)
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
				form := strings.NewReader("parent_comment_id=123&content=test")
				req := httptest.NewRequest(http.MethodPost, "/add-reply", form)
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				ctx := context.WithValue(req.Context(), "user_uuid", "test-user")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusOK, // Handler returns without error response
			expectedBody:   "",
		},
		{
			name:   "Session cookie mismatch",
			method: http.MethodPost,
			setupRequest: func() *http.Request {
				form := strings.NewReader("parent_comment_id=123&content=test")
				req := httptest.NewRequest(http.MethodPost, "/add-reply", form)
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				req.AddCookie(&http.Cookie{Name: "session_id", Value: "cookie-session"})
				ctx := context.WithValue(req.Context(), "session_id", "context-session")
				ctx = context.WithValue(ctx, "user_uuid", "test-user")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusOK, // Handler returns without error response
			expectedBody:   "",
		},
		{
			name:   "Empty content",
			method: http.MethodPost,
			setupRequest: func() *http.Request {
				form := strings.NewReader("parent_comment_id=123&content=")
				req := httptest.NewRequest(http.MethodPost, "/add-reply", form)
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
				form := strings.NewReader("parent_comment_id=123&content=test+reply")
				req := httptest.NewRequest(http.MethodPost, "/add-reply", form)
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session"})
				ctx := context.WithValue(req.Context(), "session_id", "test-session")
				ctx = context.WithValue(ctx, "user_uuid", "test-user")
				return req.WithContext(ctx)
			},
			mockAddReply: func(parentID, userID, content string) (int64, error) {
				if parentID != "123" || userID != "test-user" || content != "test reply" {
					return 0, fmt.Errorf("unexpected arguments")
				}
				return 456, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"comment_id":456,"message":"Comment added successfully"}`,
		},
		{
			name:   "Database error",
			method: http.MethodPost,
			setupRequest: func() *http.Request {
				form := strings.NewReader("parent_comment_id=123&content=test+reply")
				req := httptest.NewRequest(http.MethodPost, "/add-reply", form)
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				req.AddCookie(&http.Cookie{Name: "session_id", Value: "test-session"})
				ctx := context.WithValue(req.Context(), "session_id", "test-session")
				ctx = context.WithValue(ctx, "user_uuid", "test-user")
				return req.WithContext(ctx)
			},
			mockAddReply: func(parentID, userID, content string) (int64, error) {
				return 0, fmt.Errorf("database failure")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to add comment: database failure",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockAddReply != nil {
				models.AddReply = tt.mockAddReply
			} else {
				models.AddReply = originalAddReply
			}

			req := tt.setupRequest()
			rr := httptest.NewRecorder()

			handlers.AddReplyHandler(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			body := strings.TrimSpace(rr.Body.String())
			if tt.expectedBody != "" && body != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, body)
			}

			if tt.expectedStatus == http.StatusOK && tt.mockAddReply != nil {
				if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
					t.Errorf("expected Content-Type application/json, got %s", contentType)
				}
			}
		})
	}
}