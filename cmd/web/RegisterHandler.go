package main

import (
	"html/template"
	"net/http"

	middleware "learn.zone01kisumu.ke/git/clomollo/forum/Middleware"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
	"learn.zone01kisumu.ke/git/clomollo/forum/utils"
)

var Tmpl *template.Template

func (dep *Dependencies)RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		csrfToken := r.Context().Value("csrf_token").(string)
		Tmpl.ExecuteTemplate(w, "register.html", map[string]interface{}{
			"CSRFToken": csrfToken,
		})
		return
	}
	if r.Method != http.MethodPost {
		dep.ClientError(w, http.StatusMethodNotAllowed)
		return
	}
	if !middleware.ValidateCSRFToken(r) {
		dep.ClientError(w, http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		dep.ClientError(w, http.StatusBadRequest)
		return
	}

	// get the form data
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")

	if email == "" || username == "" || password == "" {
		dep.ClientError(w,http.StatusBadRequest)
		return
	}

	if !utils.ValidateEmail(email) {
		dep.ErrorLog.Println("Error could not validate email format")
		dep.ClientError(w,http.StatusBadRequest)
		return
	}

	userByEmail, err := models.GetUserByEmail(email)
	if err != nil {
		dep.ServerError(w, err)
		return
	}
	if userByEmail != nil {
		dep.ClientError(w,http.StatusBadRequest)
		return
	}

	userByUsername, err := models.GetUserByUsername(username)
	if err != nil {
		dep.ServerError(w, err)
		return
	}
	if userByUsername != nil {
		dep.ClientError(w,http.StatusBadRequest)
		return
	}

	if len(password) < 8 {
		dep.ClientError(w, http.StatusBadRequest)
		return
	}

	if err := models.CreateUser(email, username, password); err != nil {
		dep.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
