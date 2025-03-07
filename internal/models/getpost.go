package models

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	// "learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

// var database *models.ForumModel
var DB *sql.DB

// var f *ForumModel

func RenderPostsPage(w http.ResponseWriter,r *http.Request) {
	// if r.Method == http.MethodGet {
		var categories []struct {
			ID   string
			Name string
		}
		categoryRows, err := DB.Query("SELECT id, name FROM categories ORDER BY name")
		if err != nil {
			http.Error(w, "Failed to load categories", http.StatusInternalServerError)
			return 
		}
		defer categoryRows.Close()

		for categoryRows.Next() {
			var cat struct {
				ID   string
				Name string
			}
			if err := categoryRows.Scan(&cat.ID, &cat.Name); err != nil {
				continue
			}
			categories = append(categories, cat)
		}

		rows, err := DB.Query(`
            SELECT p.post_id, u.username, p.title, p.content, p.media, p.content_type, p.created_at 
            FROM posts p 
            JOIN users u ON p.user_uuid = u.user_uuid
        `)
		if err != nil {
			http.Error(w, "Failed to load posts", http.StatusInternalServerError)
			// return
		}
		defer rows.Close()

		var posts []map[string]interface{}
		for rows.Next() {
			var id, username, title, content string
			var createdAt time.Time
			var media []byte
			var contentType *string

			if err := rows.Scan(&id, &username, &title, &content, &media, &contentType, &createdAt); err != nil {
				fmt.Println(err)
				http.Error(w, "Failed to parse posts", http.StatusInternalServerError)
				return 
			}

			// Convert media to base64 if it exists
			// Convert media to base64 if it exists
			var mediaBase64 string
			if len(media) > 0 {
				mediaBase64 = base64.StdEncoding.EncodeToString(media)
			}

			// Handle contentType
			var contentTypeStr string
			if contentType != nil {
				contentTypeStr = *contentType // Dereference the pointer
			} else {
				contentTypeStr = "" // Or set a default value if needed
			}

			categoryRows, err := DB.Query(`
                SELECT c.id, c.name 
                FROM categories c 
                JOIN post_categories pc ON c.name = pc.category_id 
                WHERE pc.post_id = ?`, id)
			if err != nil {
				http.Error(w, "Failed to fetch post categories", http.StatusInternalServerError)
				return 
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

			var likes, dislikes int
			err = DB.QueryRow("SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND type = 'like'", id).Scan(&likes)
			if err != nil {
				http.Error(w, "Failed to fetch likes", http.StatusInternalServerError)
				return 
			}

			err = DB.QueryRow("SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND type = 'dislike'", id).Scan(&dislikes)
			if err != nil {
				http.Error(w, "Failed to fetch dislikes", http.StatusInternalServerError)
				return 
			}

			commentRows, err := GetAllCommentsForPost(id)
			if err != nil {
				http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
				return 
			}

			var comments []map[string]interface{}

			for _, comment := range commentRows {
				var commentLikes, commentDislikes int

				query := `SELECT 
					(SELECT COUNT(*) FROM comment_likes WHERE comment_id = c.id AND type = 'like') AS likes,
					(SELECT COUNT(*) FROM comment_likes WHERE comment_id = c.id AND type = 'dislike') AS dislikes
				FROM comments c
				WHERE c.id = ?;`

				err := DB.QueryRow(query, comment.ID).Scan(&commentLikes, &commentDislikes)
				if err != nil {
					http.Error(w, "Failed to parse comment likes", http.StatusInternalServerError)
					return
				}

				replies, err := CommentREplies(comment, w)
				if err != nil {
					http.Error(w, "Failed to fetch Replies", http.StatusInternalServerError)
					return
				}

				comments = append(comments, map[string]interface{}{
					"ID":              comment.ID,
					"CommentUsername": comment.UserName,
					"CommentInitial":  string(comment.UserName[0]),
					"Content":         comment.Content,
					"CreatedAt":       comment.CreatedAt,
					"Likes":           commentLikes,
					"Dislikes":        commentDislikes,
					"Replies":         replies,
				})

			}			

			posts = append(posts, map[string]interface{}{
				"ID":             id,
				"Title":          title,
				"Content":        content,
				"Likes":          likes,
				"Dislikes":       dislikes,
				"Comments":       comments,
				"CommentsLenght": len(comments),
				"Username":       username,
				"Initial":        string(username[0]),
				"Categories":     postCategories,
				"Media":          mediaBase64,
				"ContentType":    contentTypeStr,
				"CreatedAt":      createdAt,
			})
		}
        


		data := map[string]interface{}{
			"Posts":      posts,
			"Categories": categories,
		}

		userId:=r.Context().Value("user_uuid").(string)

		query:=`
		SELECT u.username
		FROM users u 
		WHERE u.user_uuid=?`
        var username string
		err=DB.QueryRow(query,userId).Scan(&username)
		if err ==nil{
           data["UserName"] = username
		   data["Initial"] = string(username[0])
		}

		RenderTemplates(w, "posts.html", data)
	}
// }

func CommentREplies(comment Comment, w http.ResponseWriter) ([]map[string]interface{}, error) {
	replyRow, err := GetAllRepliesForComment(comment.ID)
	if err != nil {
		http.Error(w, "Failed to fetch Replies", http.StatusInternalServerError)
		return nil, err
	}

	var replies []map[string]interface{}

	for _, reply := range replyRow {

		var replyLikes, replyDislikes int

		query := `SELECT 
			(SELECT COUNT(*) FROM comment_likes WHERE comment_id = r.id AND type = 'like') AS likes,
			(SELECT COUNT(*) FROM comment_likes WHERE comment_id = r.id AND type = 'dislike') AS dislikes
			FROM comments r
			WHERE r.id = ?;`

		err := DB.QueryRow(query, reply.ID).Scan(&replyLikes, &replyDislikes)
		if err != nil {
			http.Error(w, "Failed to parse reply likes", http.StatusInternalServerError)
			fmt.Println(err)
			return nil, err
		}

		replies2, err := CommentREplies(reply, w)
		if err != nil {
			http.Error(w, "Failed to fetch Replies", http.StatusInternalServerError)
			return nil, err
		}

		replies = append(replies, map[string]interface{}{
			"ID":            reply.ID,
			"ReplyUsername": reply.UserName,
			"ReplyInitial":  string(reply.UserName[0]),
			"Content":       reply.Content,
			"CreatedAt":     reply.CreatedAt,
			"Likes":         replyLikes,
			"Dislikes":      replyDislikes,
			"Replies":       replies2,
		})

	}

	return replies, nil
}





// type postCategory struct {
// 	ID   string
// 	Name string
// }

// type Post struct {
// 	Id             int
// 	PostId         string
// 	UserId         string
// 	Media          []byte
// 	Category       []string
// 	Title          string
// 	Content        string
// 	ContentType    string
// 	TimeStamp      string
// 	Likes          int
// 	Dislikes       int
// 	Comments       Comment
// 	CommentsLenght int
// 	Username       string
// 	Initial        string
// 	Categories     postCategory
// 	CreatedAt      time.Time
// }