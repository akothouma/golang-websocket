package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		fmt.Println("Session token:", cookie.Value)

		session, err := models.GetSession(cookie.Value)
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
		ctx := context.WithValue(r.Context(), "user_id", session.UserID)
		// If the session is valid, call the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CreateSession(w http.ResponseWriter, r *http.Request, userID int) {
	// Generate a new session ID
	sessionID, err := models.CreateSession(userID)
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
