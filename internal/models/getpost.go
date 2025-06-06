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
	if r.URL.Path != "/"{
		RenderTemplates(w, "error.html", map[string]string{
			"Code":"404 Page Not Found",
		})
		return
	}
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
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})

	// username := ""

	// Check if the user is logged in
	username, _ := LogedInUser(r)
	// Get user details
	f := &ForumModel{DB: DB}
	user, err := f.GetUserByUsername(username)
	if err == nil && user != nil {
		if user.ProfilePicture != "" {
			data["ProfilePicture"] = user.ProfilePicture
			data["ContentType"] = user.ContentType
		} else if len(username) > 0 {
			data["Initial"] = string(username[0])
		}
	}

	//pass csrf_token through the context
	csrfToken := r.Context().Value("csrf_token").(string)
	data["CSRFToken"] = csrfToken

	// Add debugging information
	
	data["UserName"] = username 
	data["ViewType"] = "all"
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

func Post_Categories(id string) ([]postCategory, error) {
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
			ID:   catID,
			Name: catName,
		})
	}

	return postCategories, nil
}

func RenderLikedPostsPage(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	// Check if the user is logged in
	username, _ := LogedInUser(r)
	// Get user details
	f := &ForumModel{DB: DB}
	user, err := f.GetUserByUsername(username)
	if err == nil && user != nil {
		if user.ProfilePicture != "" {
			data["ProfilePicture"] = user.ProfilePicture
			data["ContentType"] = user.ContentType
		} else if len(username) > 0 {
			data["Initial"] = string(username[0])
		}
	}

	// Get user UUID from username
	var userUUID string
	err = DB.QueryRow("SELECT user_uuid FROM users WHERE username = ?", username).Scan(&userUUID)
	if err != nil {
		http.Error(w, "Failed to get user information", http.StatusInternalServerError)
		return
	}

	// Query to get posts liked by the user
	query := `
    SELECT p.id, p.post_id, p.user_uuid, p.username, p.title, p.content, p.media, p.content_type, p.created_at
    FROM posts p
    JOIN post_likes pl ON p.post_id = pl.post_id
    JOIN users u ON p.user_uuid = u.user_uuid
    WHERE pl.user_id = ? AND pl.type = 'like'
    ORDER BY p.created_at DESC
`

	rows, err := DB.Query(query, userUUID)
	if err != nil {
		http.Error(w, "Failed to fetch liked posts", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	defer rows.Close()

    likedPosts, err := PostsRows(rows)
    if err != nil {
        http.Error(w, "Failed to fetch liked posts "+err.Error(), http.StatusInternalServerError)
        fmt.Println(err)
        return
    }
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

	data["UserName"] = username
	data["Initial"] = string(username[0])
	data["ViewType"] = "liked"
	data["Posts"] = likedPosts
	data["Categories"] = categories

	RenderTemplates(w, "index.html", data)
}

func RenderMyPostsPage(w http.ResponseWriter, r *http.Request) {
	
	data := make(map[string]interface{})

	// Check if the user is logged in
	username, _ := LogedInUser(r)
	// Get user details
	f := &ForumModel{DB: DB}
	user, err := f.GetUserByUsername(username)
	if err == nil && user != nil {
		if user.ProfilePicture != "" {
			data["ProfilePicture"] = user.ProfilePicture
			data["ContentType"] = user.ContentType
		} else if len(username) > 0 {
			data["Initial"] = string(username[0])
		}
	}

	// Get user UUID from username
	var userUUID string
	err = DB.QueryRow("SELECT user_uuid FROM users WHERE username = ?", username).Scan(&userUUID)
	if err != nil {
		http.Error(w, "Failed to get user information", http.StatusInternalServerError)
		return
	}

	// Query to get posts created by the user
	query := `
        SELECT p.id, p.post_id, p.user_uuid, p.username, p.title, p.content, p.media, p.content_type, p.created_at
        FROM posts p
        JOIN users u ON p.user_uuid = u.user_uuid
        WHERE p.user_uuid = ?
        ORDER BY p.created_at DESC
    `

	rows, err := DB.Query(query, userUUID)
	if err != nil {
		http.Error(w, "Failed to fetch your posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

    myPosts, err := PostsRows(rows)
    if err != nil {
        http.Error(w, "Failed to fetch your posts PostRows"+err.Error(), http.StatusInternalServerError)
        fmt.Println(err)
        return
    }
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
	

	data["UserName"] = username
	data["Initial"] = string(username[0])
	data["ViewType"] = "mine"
	data["Posts"] = myPosts
	data["Categories"] = categories

	RenderTemplates(w, "index.html", data)
}
