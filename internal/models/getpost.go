package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var DB *sql.DB

type AllData struct {
	Post       []PostData
	Categories []PostCategory
	Comments []Comments
}
type PostCategory struct {
	ID   string
	Name string
}

type PostData struct {
	Id          string
	Username    string
	Title       string
	Content     string
	CreatedAt   time.Time
	ContentType *string
	Likes int
	Dislikes int
}
type Comments struct {
	CommentId, 
	CommentContent string
	CommentCreatedAt time.Time
	CommentLikes int
	CommentDislikes int
}

func RenderPostsPage() (AllData, error) {
	var allData AllData
	categoryRows, err := DB.Query("SELECT id, name FROM categories ORDER BY name")
	if err != nil {
		// Error("Failed to load categories", http.StatusInternalServerError)
		return allData, errors.New("Failed to load categories")
	}
	defer categoryRows.Close()

	allCategories := []PostCategory{} //Populate to all data
	for categoryRows.Next() {
		var cat PostCategory

		if err := categoryRows.Scan(&cat.ID, &cat.Name); err != nil {
			continue
		}
		allCategories = append(allCategories, cat)
	}

	rows, err := DB.Query(`
            SELECT p.post_id, u.username, p.title, p.content, p.content_type, p.created_at 
            FROM posts p 
            JOIN users u ON p.user_uuid = u.user_uuid
        `)
	if err != nil {
		// http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return allData, errors.New("Failed to load posts")
	}
	defer rows.Close()

	var posts PostData
	var AllPosts []PostData
	var AllComments []Comments
	for rows.Next() {
		if err := rows.Scan(&posts.Id, &posts.Username, &posts.Title, &posts.Content, &posts.ContentType, &posts.CreatedAt); err != nil {
			// fmt.Println(err)
			// http.Error(w, "Failed to parse posts", http.StatusInternalServerError)
			return allData, fmt.Errorf("Failed to parse posts:\n%v", err)
		}

		// Handle contentType
		// var contentTypeStr string
		// if posts.ContentType != nil {
		// 	contentTypeStr = *posts.ContentType // Dereference the pointer
		// } else {
		// 	contentTypeStr = "" // Or set a default value if needed
		// }
		 

		categoryRows, err := DB.Query(`
                SELECT c.id, c.name 
                FROM categories c 
                JOIN post_categories pc ON c.name = pc.category_id 
                WHERE pc.post_id = ?`, posts.Id)
		if err != nil {
			// http.Error(w, "Failed to fetch post categories", http.StatusInternalServerError)
			return allData, errors.New("Failed to fetch post categories")
		}
		defer categoryRows.Close()

		var postCategories []map[string]string
		for categoryRows.Next() {
			var catID, catName string
			if err := categoryRows.Scan(&catID, &catName); err != nil {
				continue
			}
			postCategories = append(postCategories, map[string]string{
				"ID":   catID,
				"Name": catName,
			})
		}

		
		err = DB.QueryRow("SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND type = 'like'", posts.Id).Scan(&posts.Likes)
		if err != nil {
			// http.Error(w, "Failed to fetch likes", http.StatusInternalServerError)
			// return
			return allData, errors.New("Failed to fetch likes")
		}

		err = DB.QueryRow("SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND type = 'dislike'", posts.Id).Scan(&posts.Dislikes)
		if err != nil {
			// http.Error(w, "Failed to fetch dislikes", http.StatusInternalServerError)
			// return
			return allData, errors.New("Failed to fetch dislikes")
		}

		commentRows, err := DB.Query(`
                SELECT c.id, c.content, c.created_at,
                (SELECT COUNT(*) FROM comment_likes WHERE comment_id = c.id AND type = 'like') as likes,
                (SELECT COUNT(*) FROM comment_likes WHERE comment_id = c.id AND type = 'dislike') as dislikes
                FROM comments c
                WHERE c.post_id = ?
                ORDER BY c.created_at DESC`, posts.Id)
		if err != nil {
			// http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
			// return
			return allData, errors.New("Failed to fetch comments")
		}
		defer commentRows.Close()

		var comment Comments
		for commentRows.Next() {
			
			if err := commentRows.Scan(&comment.CommentId, &comment.CommentContent, &comment.CommentCreatedAt, &comment.CommentLikes, &comment.CommentDislikes); err != nil {
				// http.Error(w, "Failed to parse comment", http.StatusInternalServerError)
				return allData, errors.New("Failed to parse comments")
			}
			AllComments = append(AllComments, comment)
		}
		//Populate all posts with post data
		AllPosts = append(AllPosts, posts)

	

		// posts = append(posts, map[string]interface{}{
		// 	"ID":          id,
		// 	"Title":       title,
		// 	"Content":     content,
		// 	"Likes":       likes,
		// 	"Dislikes":    dislikes,
		// 	"Comments":    comments,
		// 	"Username":    username,
		// 	"Categories":  postCategories,
		// 	"ContentType": contentTypeStr,
		// 	"CreatedAt":   createdAt,
		// })
	}
	
	allData.Categories = allCategories
	allData.Post = AllPosts
	allData.Comments = AllComments
	// var data map[string]interface{}
	// data = map[string]interface{}{
	// 	"Posts":      posts,
	// 	"Categories": categories,
	// }

	// // RenderTemplates(w, "posts.html", data)
	// postHtml, err := template.ParseFiles("./ui/html/posts.html")
	// if err != nil {
	// 	// fmt.Println("Error loading posts html\n",err )
	// 	fmt.Println("")
	// }
	// postHtml.Execute(w, data)
	// }
	return allData, nil
}
