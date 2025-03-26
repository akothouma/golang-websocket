package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/google/uuid"
	"learn.zone01kisumu.ke/git/clomollo/forum/utils"
)

type RegisterRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Tac             bool   `json:"tac"`
	Csrf            string `json:"csrfToken"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (dep *Dependencies) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var regReq RegisterRequest
	registerTemp, err := template.ParseFiles("./ui/html/register.html")
	if err == nil {
		if r.Method == http.MethodGet {
		csrfToken := r.Context().Value("csrf_token").(string)
			registerTemp.ExecuteTemplate(w, "register.html", map[string]interface{}{
				"CSRFToken":csrfToken ,
			})
			fmt.Println("here",csrfToken)
			return
		}
		if r.Method != http.MethodPost {
			dep.ClientError(w, http.StatusMethodNotAllowed)
			return
		}
		
		if err := r.ParseForm(); err != nil {
			dep.ClientError(w,http.StatusInternalServerError)
			return
		}
		
		if err := json.NewDecoder(r.Body).Decode(&regReq); err != nil {
			dep.ClientError(w, http.StatusBadRequest)
		}
		
		if !dep.ValidateCSRFToken(r,regReq.Csrf) {
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid CSRF token"})
			return
		}

		userUuid := uuid.New().String()
		email := regReq.Email
		username := regReq.Username
		password := regReq.Password
		confirmPassword := regReq.ConfirmPassword
		tac := regReq.Tac

		fmt.Println(userUuid, email, username, password)
		if email == "" || username == "" || password == "" || confirmPassword == "" {
			json.NewEncoder(w).Encode(ErrorResponse{Error: "All fields are required"})
			return
		}
		if !tac {
			dep.ClientError(w, http.StatusBadRequest)
			return
		}

		if password != confirmPassword {
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Passwords do not match"})
			return
		}

		if !utils.ValidateEmail(email) {
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Please provide a valid email address"})
			return
		}

		userByEmail, err := dep.Forum.GetUserByEmail(email)
		if err != nil {
			dep.ServerError(w, err)
			return
		}
		if userByEmail != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Email already exists"})
			return
		}

		userByUsername, err := dep.Forum.GetUserByUsername(username)
		if err != nil {
			dep.ServerError(w, err)
			return
		}
		if userByUsername != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Username already taken"})
			return
		}

		if len(password) < 8 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Password must be at least 8 characters"})
			return
		}
		if err := dep.Forum.CreateUser(userUuid, email, username, password); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to create user"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message":  "Registration succesful.Redirecting to login...",
			"redirect": "/login",
		})
	}

	// http.Redirect(w, r, "/login", http.StatusSeeOther)
}
