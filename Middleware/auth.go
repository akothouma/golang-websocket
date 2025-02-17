package middleware

import (
	"context"
	"fmt"
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
