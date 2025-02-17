package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func CSRFMiddleware(next http.Handler) http.Handler {
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


func ValidateCSRFToken(r *http.Request) bool {
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