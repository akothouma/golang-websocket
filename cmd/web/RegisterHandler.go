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
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
    ConfirmPassword string `json:"confirmPassword"`
	Tac bool `json:"tac"`
	Csrf *http.Request `json:"csrfToken"`
}


type ErrorResponse struct {
    Error string `json:"error"`
}

func (dep *Dependencies) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	registerTemp, err := template.ParseFiles("./ui/html/register.html")
	if err == nil {
		if r.Method == http.MethodGet {
			csrfToken := r.Context().Value("csrf_token").(string)
			registerTemp.ExecuteTemplate(w, "register.html", map[string]interface{}{
	 			"CSRFToken": csrfToken,
			})
	 		return
		}
	var regReq RegisterRequest
		if r.Method != http.MethodPost {
			dep.ClientError(w, http.StatusText(http.StatusMethodNotAllowed),http.StatusMethodNotAllowed)
			return
		}
		if !dep.ValidateCSRFToken(regReq.Csrf) {
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid CSRF token"})
			return
		}

		// if err := r.ParseForm(); err != nil {
		// 	dep.ClientError(w, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)
		// 	return
		// }


		if err:=json.NewDecoder(r.Body).Decode(&regReq);err!=nil{
			dep.ClientError(w,http.StatusText(http.StatusBadRequest),http.StatusBadRequest)
		}

		userUuid := uuid.New().String()
		email := regReq.Email
		username := regReq.Username
		password := regReq.Password
		confirmPassword:=regReq.ConfirmPassword
		tac:=regReq.Tac
		
        fmt.Println(userUuid,email,username,password)
		if email == "" || username == "" || password == ""  || confirmPassword==""{
			json.NewEncoder(w).Encode(ErrorResponse{Error: "All fields are required"})
			return
		}
		if !tac{
            dep.ClientError(w,"Accept terms and conditions",http.StatusBadRequest)
			return
		}

		if password != confirmPassword{
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
			dep.ClientError(w, "Invalid credentials",http.StatusBadRequest)
			return
		}


		userByUsername, err := dep.Forum.GetUserByUsername(username)
		if err != nil {
			dep.ServerError(w, err)
			return
		}
		if userByUsername != nil {
			dep.ClientError(w,"Invalid credentials", http.StatusBadRequest)
			return
		}

		if len(password) < 8 {
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Password must be at least 8 characters"})
			return
		}
		if err := dep.Forum.CreateUser(userUuid,email, username, password); err != nil {

			dep.ServerError(w, err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message":"Registration succesful",
			"redirect":"/login",
		})
	}

	//http.Redirect(w, r, "/login", http.StatusSeeOther)

}