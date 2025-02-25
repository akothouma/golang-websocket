package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (dep *Dependencies) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		fmt.Println("Session token:", cookie.Value)

		session, err := dep.Forum.GetSession(cookie.Value)
		if err != nil || session.ExpiresAt.Before(time.Now()) {
			http.SetCookie(w, &http.Cookie{
				Name:    "session_id",
				Value:   "",
				Path:    "/",
				Expires: time.Unix(0, 0),
			})
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		// Add the user ID to the request context
		ctx := context.WithValue(r.Context(), "user_uuid", session.UserID)
		ctx = context.WithValue(ctx, "session_id", session.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
		
	})
}

func (dep *Dependencies) CreateSession(w http.ResponseWriter, r *http.Request, userID string) {
	// Generate a new session ID
	sessionID, err := dep.Forum.CreateSession(userID)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Set the session ID as a cookie
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)
}

func (dep *Dependencies) CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("csrf_token")
		var csrfToken string

		if err != nil || cookie.Value == "" {
			csrfToken = uuid.New().String()

			http.SetCookie(w, &http.Cookie{
				Name:     "csrf_token",
				Value:    csrfToken,
				Path:     "/",
				HttpOnly: true,
				Secure:   false,
				Expires:  time.Now().Add(24 * time.Hour),
			})
		} else {
			csrfToken = cookie.Value
		}

		// Add the CSRF token to the request context
		ctx := context.WithValue(r.Context(), "csrf_token", csrfToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (dep *Dependencies) ValidateCSRFToken(r *http.Request) bool {
	// ValidateCSRFToken checks if the CSRF token is valid
	// Get the CSRF token from the form
	formToken := r.FormValue("csrf_token")
	log.Printf("CSRF token from form: %s\n", formToken)
	if formToken == "" {
		return false
	}

	// Get the CSRF token from the cookie
	cookie, err := r.Cookie("csrf_token")
	if err != nil || cookie.Value == "" {
		return false
	}

	// Compare the tokens
	return formToken == cookie.Value
}
