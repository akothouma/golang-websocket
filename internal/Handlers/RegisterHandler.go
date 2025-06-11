package handlers

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
	"learn.zone01kisumu.ke/git/clomollo/forum/utils"
)

type RegisterRequest struct {
    FirstName   string `json:"firstName"`
    LastName    string `json:"lastName"`
    Username    string `json:"username"`
    Email       string `json:"email"`
    Age         int    `json:"age,omitempty"`    // omitempty because client might send null if age is not entered,
                                               // and it allows the field to be absent in JSON too.
    Gender      string `json:"gender,omitempty"` // similar to age, allows absence or null to become ""
    Password    string `json:"password"`
    Tac         bool   `json:"tac"`
    CsrfToken   string `json:"csrfToken"`      // Renamed struct field for clarity, tag matches client
}
type ErrorResponse struct {
	Error string `json:"error"`
}

type Dependencies struct {
	ErrorLog  *log.Logger
	InfoLog   *log.Logger
	Forum     *models.ForumModel
	Templates *template.Template
	db        *sql.DB
}

func (dep *Dependencies) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// // The GET part for serving the HTML page can remain similar if you still have a separate register.html
	// // However, for a pure SPA driven by JS fetch, GET might not be needed or would serve the main index.html.
	// // For this example, I'll assume POST is the primary concern for registration API.
	// if r.Method == http.MethodGet {
	// 	// If you're using Go templates for SPA shell or individual views, handle it here.
	// 	// For a pure API backend, GET on /register might not do anything or might return an error.
	// 	// For simplicity, let's assume SPA, so GET /register may not be a standard operation.
	// 	http.ServeFile(w,r, "./ui/html/index.html") // Or however you serve your SPA entry point
	// 	return
	// }\
// if r.Method == http.MethodGet {
// 		csrfToken := r.Context().Value("csrf_token").(string)
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(map[string]string{
// 			"csrfToken": csrfToken,
// 		})
		
// 		return
// 	}


	if r.Method != http.MethodPost {
		dep.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	var regReq RegisterRequest
	// Decode JSON from request body
	err := json.NewDecoder(r.Body).Decode(&regReq)
	if err != nil {
		dep.ErrorLog.Printf("Failed to decode registration request: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request format. Check JSON structure."})
		return
	}
	// It's good practice to close the body if not already handled by a framework
	defer r.Body.Close()


	// Log the received request data for debugging
	dep.InfoLog.Printf("Received registration request: %+v", regReq)


	// CSRF Token Validation
	// Assuming your middleware places the server-side token in the request context.
	// And your client sends it in the JSON body as "csrfToken".
	if !dep.ValidateCSRFToken(r, regReq.CsrfToken) { // regReq.Csrf comes from `json:"csrfToken"`
		dep.ErrorLog.Println("Invalid CSRF token during registration")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden) // 403 Forbidden is appropriate for CSRF failure
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid security token. Please refresh and try again."})
		return
	}

	// Basic Validations (Password confirmation is client-side, but can be re-checked)
	// You removed ConfirmPassword from struct, so this check needs plainPassword & confirmPassword from form if done server-side
	// For now, relying on client-side check.

	// Field presence validation (basic) - Add FirstName and LastName if they are mandatory
	if regReq.Email == "" || regReq.Username == "" || regReq.Password == "" || regReq.FirstName == "" || regReq.LastName == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Required fields (First Name, Last Name, Email, Username, Password) cannot be empty."})
		return
	}

	if !regReq.Tac {
		dep.InfoLog.Printf("Registration attempt without accepting TAC: User %s", regReq.Username)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest) // Bad Request
		json.NewEncoder(w).Encode(ErrorResponse{Error: "You must accept the terms and conditions to register."})
		return
	}

	if !utils.ValidateEmail(regReq.Email) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Please provide a valid email address."})
		return
	}

	// Check if email already exists
	existingUserByEmail, err := dep.Forum.GetUserByEmail(regReq.Email)
	if err != nil && err != sql.ErrNoRows { // sql.ErrNoRows is expected if user doesn't exist
		dep.ServerError(w, fmt.Errorf("error checking email existence: %w", err))
		return
	}
	if existingUserByEmail != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict) // 409 Conflict is better for existing resource
		json.NewEncoder(w).Encode(ErrorResponse{Error: "This email address is already registered."})
		return
	}

	// Check if username already exists
	existingUserByUsername, err := dep.Forum.GetUserByUsername(regReq.Username)
	if err != nil && err != sql.ErrNoRows {
		dep.ServerError(w, fmt.Errorf("error checking username existence: %w", err))
		return
	}
	if existingUserByUsername != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict) // 409 Conflict
		json.NewEncoder(w).Encode(ErrorResponse{Error: "This username is already taken. Please choose another."})
		return
	}

	// Password Policy
	if len(regReq.Password) < 8 { // Example: minimum length
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Password must be at least 8 characters long."})
		return
	}

	// Generate a new UUID for the user
	userUuid := uuid.New().String()

	// Call the updated CreateUser function with all the new fields
	err = dep.Forum.CreateUser(
		userUuid,
		regReq.FirstName, // New
		regReq.LastName,  // New
		regReq.Email,
		regReq.Username,
		regReq.Password,  // This is the plain password, CreateUser will hash it
		regReq.Age,       // New (will be 0 if not provided and not omitempty and no default in DB)
		regReq.Gender,    // New (will be "" if not provided)
	)

	if err != nil {
		dep.ErrorLog.Printf("Failed to create user in database (%s): %v", regReq.Username, err)
		// Check for specific database errors if your CreateUser doesn't already refine them
		// (e.g., SQLite UNIQUE constraint errors often contain "UNIQUE constraint failed")
		if strings.Contains(strings.ToLower(err.Error()), "unique constraint failed") {
			// This might be redundant if GetUserByEmail/Username checks are solid,
			// but good for race conditions or other DB-level uniqueness.
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "An account with this email or username already exists."})
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError) // 500 for other DB errors
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Could not create account at this time. Please try again later."})
		}
		return
	}

	dep.InfoLog.Printf("User successfully registered: UUID %s, Username %s", userUuid, regReq.Username)

	// Send success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created is conventional for successful resource creation
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Registration successful! You can now login.",
		"redirect": "/login", // Or a path that shows the login form in your SPA
	})
}