package handlers

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

func PostHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodGet {
		csrfToken := r.Context().Value("csrf_token").(string)
		Tmpl.ExecuteTemplate(w, "login.html", map[string]interface{}{
			"CSRFToken": csrfToken,
		})
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return 
	}
	postContent := r.FormValue("postContent")
	postId := uuid.New().String()
	category := r.Form["categories"]
	title := r.FormValue("title")
	userId := r.Context().Value("user_id")

	if userId != nil {
		models.CreatePost(postId, title, postContent, category)
	}
	fmt.Println("You have to be signed in to be able to post")
}
