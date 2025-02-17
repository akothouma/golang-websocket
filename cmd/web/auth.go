package main

import (
	"html/template"
	"log"
	"net/http"

	middleware "learn.zone01kisumu.ke/git/clomollo/forum/Middleware"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
	"learn.zone01kisumu.ke/git/clomollo/forum/utils"
)

var tmpl = template.Must(template.ParseGlob("ui/html/*.html"))

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		csrfToken := r.Context().Value("csrf_token").(string)
		tmpl.ExecuteTemplate(w, "register.html", map[string]interface{}{
			"CSRFToken": csrfToken,
		})
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !middleware.ValidateCSRFToken(r) {
		http.Error(w, "Invalid CSRF token", http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Cannot parse form", http.StatusBadRequest)
		return
	}

	// get the form data
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")

	if email == "" || username == "" || password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	if !utils.ValidateEmail(email) {
		log.Println("Error could not validate email format")
		http.Error(w, "Invalid email address", http.StatusBadRequest)
		return
	}

	userByEmail, err := models.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if userByEmail != nil {
		http.Error(w, "Email already taken", http.StatusBadRequest)
		return
	}

	userByUsername, err := models.GetUserByUsername(username)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if userByUsername != nil {
		http.Error(w, "Username already taken", http.StatusBadRequest)
		return
	}

	if len(password) < 8 {
		http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	if err := models.CreateUser(email, username, password); err != nil {
		http.Error(w, "Cannot create user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}


func LoginHandler(w http.ResponseWriter, r *http.Request){
	if r.Method==http.MethodGet{
		csrfToken:=r.Context().Value("csrf_token").(string)
		tmpl.ExecuteTemplate(w,"login.html",map[string]interface{}{
			"CSRFToken":csrfToken,
		})
		return
	}
	if r.Method != http.MethodPost{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !middleware.ValidateCSRFToken(r){
		http.Error(w,"Invalid CSRF token", http.StatusForbidden)
		return
	}

	if err:=r.ParseForm(); err!=nil{
		http.Error(w, "Cannot parse form", http.StatusBadRequest)
		return
	}

	email:=r.FormValue("email")
	password:=r.FormValue("password")

	if email=="" || password==""{
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	user,err:=models.GetUserByEmail(email)
	if err!=nil{
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if user==nil{
		http.Error(w,"Invalid credentials",http.StatusUnauthorized)
		return
	}

	if !user.CheckPassword(password){
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	middleware.CreateSession(w,r,user.ID)
	
	http.Redirect(w,r,"/",http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request){
	cookie,err:=r.Cookie("session_id")
	if err!=nil || cookie.Value==""{
		http.Redirect(w,r,"/",http.StatusSeeOther)
		return
	}
	err=models.DeleteSession(cookie.Value)
	if err!=nil{
		http.Error(w,"Failed to delete session",http.StatusInternalServerError)
		return
	}
	http.SetCookie(w,&http.Cookie{
		Name:"session_id",
		Value:"",
		Path: "/",
		Expires: time.Unix(0,0), 
	})
	http.Redirect(w,r,"/",http.StatusSeeOther)
}
