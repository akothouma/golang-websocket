package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

// /home/lakoth/forum-1/ui/html/login.html
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Csrf     string `json:"csrfToken"`
}

type ErrorResponseLogin struct {
	Error string `json:"error"`
}

func (dep *Dependencies) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var logReq LoginRequest
	loginTemplate, err := template.ParseFiles("./ui/html/index.html")
	if err != nil {
		http.Error(w, "NOT FOUND\nLogin template not found", http.StatusNotFound)
		return
	}
	if r.Method == http.MethodGet {
		csrfToken := r.Context().Value("csrf_token").(string)
		loginTemplate.ExecuteTemplate(w, "index.html", map[string]interface{}{
			"CSRFToken": csrfToken,
		})
		return
	}
	if r.Method != http.MethodPost {
		dep.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	// if err := r.ParseForm(); err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	json.NewEncoder(w).Encode(ErrorResponseLogin{Error: "Something went wrong.Try again later"})
	// 	return
	// }
	if err := json.NewDecoder(r.Body).Decode(&logReq); err != nil {
		dep.ClientError(w, http.StatusInternalServerError)
		return
	}

	if !dep.ValidateCSRFToken(r, logReq.Csrf) {
		http.Error(w, "Invalid CSRF token", http.StatusForbidden)
		return
	}
	email := logReq.Email
	password := logReq.Password

	if email == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponseLogin{Error: "All fields are required"})
		return
	}

	user, err := dep.Forum.GetUserByEmail(email)
	if err != nil {
		dep.ServerError(w, err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponseLogin{Error: "You have to register first"})
		return

	}

	if user == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponseLogin{Error: "You have to register first"})
		return
	}

	if !user.CheckPassword(password) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponseLogin{Error: "Invalid email or password"})
		return
	}
	dep.CreateSession(w, r, user.UserID)

// For AJAX requests, return success response instead of redirect
    fmt.Println("Login successful!!!")
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Login successful",
        "redirect": "/",
    })
}
