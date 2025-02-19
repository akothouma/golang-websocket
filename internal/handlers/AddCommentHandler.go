package handlers

import (
	"net/http"

	middleware "learn.zone01kisumu.ke/git/clomollo/forum/Middleware"
)

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost{
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
	}

	csrfToken := r.Header.Get("X-CSRF-Token")

	if !middleware.ValidateCSRFToken(r, csrfToken){
		http.Error(w, "Invalid CSRF token", http.StatusForbidden)
		return
	}
}
