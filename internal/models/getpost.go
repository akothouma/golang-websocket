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

	// Add debugging information
	fmt.Printf("Rendering homepage - User logged in: %v, Username: %s\n", err == nil, username)
	data["UserName"] = username 
	data["Posts"] = posts
	data["Categories"] = categories

	// fmt.Println("categories:", categories)

	RenderTemplates(w, "index.html", data)
}

func RenderProfile(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	f := &ForumModel{DB: DB}
	username, err := LogedInUser(r)
	if err != nil {
		// User not logged in, set as guest
		RenderTemplates(w, "index.html", data)
		return
	}

	// Fetch user details from the database
	user, err := f.GetUserByUsername(username)
	if err != nil {
		fmt.Println("Error fetching user details:", err)
		http.Error(w, "Failed to retrieve user details", http.StatusInternalServerError)
		return
	}

	// Populate all required template fields
	data["UserName"] = username

	// If user exists, add their details
	if user != nil {
		// Add profile picture if it exists
		if user.ProfilePicture != "" {
			data["ProfilePicture"] = user.ProfilePicture
			data["ContentType"] = user.ContentType // Make sure this field exists in your User struct
		} else {
			// Set initial for avatar if no profile pic
			if len(username) > 0 {
				data["Initial"] = string(username[0])
			} else {
				data["Initial"] = "U"
			}
		}
	}

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
			WHERE pc.postId = ?`, id)
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
			WHERE pc.postId = ?`, id)
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
	username, err := LogedInUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
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

	var likedPosts []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.PostId, &p.UserId, &p.UserName, &p.Title, &p.Content, &p.Media, &p.ContentType, &p.CreatedAt)
		if err != nil {
			http.Error(w, "Failed to fetch liked posts", http.StatusInternalServerError)
			return
		}

		p.Initial = string(p.UserName[0])

		p.MediaString = MediaToBase64(p.Media)

		p.Comments, err = GetAllCommentsForPost(p.PostId)
		if err != nil {
			http.Error(w, "Failed to fetch liked posts", http.StatusInternalServerError)
			return
		}

		p.Likes, p.Dislikes, err = PostLikesDislikes(p.PostId)
		if err != nil {
			http.Error(w, "Failed to fetch liked posts", http.StatusInternalServerError)
			return
		}

		p.Categories, err = Post_Categories(p.PostId)
		if err != nil {
			http.Error(w, "Failed to fetch liked posts", http.StatusInternalServerError)
			return
		}

		likedPosts = append(likedPosts, p)
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

	data := make(map[string]interface{})

	data["UserName"] = username
	data["Initial"] = string(username[0])

	data["Posts"] = likedPosts
	data["Categories"] = categories

	RenderTemplates(w, "index.html", data)
}

func RenderMyPostsPage(w http.ResponseWriter, r *http.Request) {
	username, err := LogedInUser(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
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

	var myPosts []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.PostId, &p.UserId, &p.UserName, &p.Title, &p.Content, &p.Media, &p.ContentType, &p.CreatedAt)
		if err != nil {
			http.Error(w, "Failed to fetch your posts", http.StatusInternalServerError)
			return
		}

		p.Initial = string(p.UserName[0])

		p.MediaString = MediaToBase64(p.Media)

		p.Comments, err = GetAllCommentsForPost(p.PostId)
		if err != nil {
			http.Error(w, "Failed to fetch your posts", http.StatusInternalServerError)
			return
		}

		p.Likes, p.Dislikes, err = PostLikesDislikes(p.PostId)
		if err != nil {
			http.Error(w, "Failed to fetch your posts", http.StatusInternalServerError)
			return
		}

		p.Categories, err = Post_Categories(p.PostId)
		if err != nil {
			http.Error(w, "Failed to fetch your posts", http.StatusInternalServerError)
			return
		}

		myPosts = append(myPosts, p)
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

	data := make(map[string]interface{})

	data["UserName"] = username
	data["Initial"] = string(username[0])

	data["Posts"] = myPosts
	data["Categories"] = categories

	RenderTemplates(w, "index.html", data)
}
