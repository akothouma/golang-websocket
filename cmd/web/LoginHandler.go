package main

import (
	"net/http"
)

func (dep *Dependencies) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		csrfToken := r.Context().Value("csrf_token").(string)
		Tmpl.ExecuteTemplate(w, "login.html", map[string]interface{}{
			"CSRFToken": csrfToken,
		})
		return
	}
	if r.Method != http.MethodPost {
		dep.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	if dep.ValidateCSRFToken(r) {
		http.Error(w, "Invalid CSRF token", http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		dep.ClientError(w, http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		dep.ClientError(w, http.StatusBadRequest)
		return
	}

	user, err := dep.Forum.GetUserByEmail(email)
	if err != nil {
		dep.ServerError(w, err)
		return
	}

	if user == nil {
		dep.ClientError(w, http.StatusUnauthorized)
		return
	}

	if !user.CheckPassword(password) {
		dep.ClientError(w, http.StatusUnauthorized)
		return
	}
	dep.CreateSession(w, r, user.ID)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
