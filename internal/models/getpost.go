package models

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	// "time"
)

// var database *models.ForumModel
var DB *sql.DB

// var f *ForumModel

func RenderPostsPage(w http.ResponseWriter, r *http.Request) {
	// if r.Method == http.MethodGet {

	var categories []postCategory
	categoryRows, err := DB.Query("SELECT id, name FROM categories ORDER BY name")
	if err != nil {
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}
	defer categoryRows.Close()

	for categoryRows.Next() {
		var cat postCategory
		if err := categoryRows.Scan(&cat.ID, &cat.Name); err != nil {
			continue
		}
		categories = append(categories, cat)
	}

	posts, err := AllPosts()
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	
	data := make(map[string]interface{})

	username := ""

	if username, err = LogedInUser(r); err != nil {
		fmt.Println(err)
	} else {
		data["UserName"] = username
		data["Initial"] = string(username[0])
	}

	data["Posts"] = posts
	data["Categories"] = categories

	// fmt.Println("categories:", categories)

	RenderTemplates(w, "index.html", data)
}

func LogedInUser(r *http.Request) (string, error) {
	session, err := r.Cookie("session_id")
	if err != nil {
		return "", fmt.Errorf("Session cookie not found")
	}

	sessionID := session.Value

	// Query to check if session is valid and fetch the username
	query := `
		SELECT u.username 
		FROM users u
		JOIN sessions s ON u.user_uuid = s.user_uuid
		WHERE s.id = ? AND s.expires_at > CURRENT_TIMESTAMP`

	var username string

	err = DB.QueryRow(query, sessionID).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No valid session found for session ID: %s, err: %v", sessionID, err)
			return "", fmt.Errorf("no valid session found for session ID: %s", sessionID)
		} else {
			log.Printf("Database error: %v", err)
			return "", fmt.Errorf("database error: %w", err)
		}
	}

	return username, nil
}

func PostLikesDislikes(id string) (int, int, error) {
	var likes, dislikes int
	err := DB.QueryRow("SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND type = 'like'", id).Scan(&likes)
	if err != nil {
		return 0, 0, fmt.Errorf("Failed to fetch likes %w", err)
	}

	err = DB.QueryRow("SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND type = 'dislike'", id).Scan(&dislikes)
	if err != nil {
		return 0, 0, fmt.Errorf("Failed to fetch dislikes %w", err)
	}

	return likes, dislikes, nil
}

func (postCategories *postCategory) AllCategories(id string) error {
	categoryRows, err := DB.Query(`
			SELECT c.id, c.name 
			FROM categories c 
			JOIN post_categories pc ON c.name = pc.category_id 
			WHERE pc.post_id = ?`, id)
	if err != nil {
		return fmt.Errorf("Failed to fetch post categories, %w", err)
	}
	defer categoryRows.Close()

	for categoryRows.Next() {
		var catID, catName string
		if err := categoryRows.Scan(&catID, &catName); err != nil {
			continue
		}
		postCategories.ID = catID
		postCategories.Name = catName

	}
	return nil
}


func Post_Categories(id string)([]postCategory, error){

	categoryRows, err := DB.Query(`
			SELECT c.id, c.name 
			FROM categories c 
			JOIN post_categories pc ON c.name = pc.category_id 
			WHERE pc.post_id = ?`, id)
	if err != nil {		
		return nil, fmt.Errorf("Failed to fetch post categories %w", err)
	}
	defer categoryRows.Close()

	var postCategories []postCategory
	for categoryRows.Next() {
		var catID, catName string
		if err := categoryRows.Scan(&catID, &catName); err != nil {
			continue
		}
		postCategories = append(postCategories, postCategory{
			ID:  catID,
			Name: catName,
		})
	}

	return postCategories, nil

}