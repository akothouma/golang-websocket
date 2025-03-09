package models

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"time"
)

//	type Post struct {
//		Id          int
//		PostId      string
//		UserId      string
//		Media       []byte
//		Category    []string
//		Title       string
//		Content string
//		ContentType string
//		TimeStamp   string
//	}
type postCategory struct {
	ID   string
	Name string
}

type Post struct {
	ID             int
	PostId         string
	UserId         string
	Media          []byte
	MediaString    string
	Category       []string
	Title          string
	Content        string
	ContentType    string
	TimeStamp      string
	Likes          int
	Dislikes       int
	Comments       []Comment
	CommentsLenght int
	UserName       string
	Initial        string
	Categories     postCategory
	CreatedAt      time.Time
}

func CreatePost(p *Post) error {
	var username string
	err := DB.QueryRow("SELECT username FROM users WHERE user_uuid = ?", p.UserId).Scan(&username)
	if err != nil {
		log.Printf("Failed to fetch username: %v", err)
		return fmt.Errorf("failed to fetch username: %w", err)
	}

	query := "INSERT INTO posts(post_id, user_uuid, username, title, content,  media, content_type) VALUES(?, ?, ?, ?, ?, ?, ?)"
	_, err = DB.Exec(query, p.PostId, p.UserId, username, p.Title, p.Content, p.Media, p.ContentType)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to insert a post")
	}

	// Insert categories
	fmt.Println("categories", p.Category)

	for _, categoryNames := range p.Category {
		_, err = DB.Exec(`
            INSERT INTO post_categories (post_id, category_id)
            VALUES (?, ?)`,
			p.PostId, categoryNames,
		)
		if err != nil {
			fmt.Println(err)

			return err
		}
	}
	fmt.Println("Your post has been succesfully created")
	return nil
}

func FindPostById(id string) (*Post, error) {
	query := "SELECT id,user_uuid,category,title,content FROM posts WHERE id=?"
	row := DB.QueryRow(query, id)
	post := Post{}
	err := row.Scan(&post.PostId, &post.UserId, &post.Title, &post.Content)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &post, nil
}

func AllPosts() ([]Post, error) {
	query := "SELECT * FROM posts"
	rows, err := DB.Query(query)
	posts := []Post{}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.PostId, &p.UserId, &p.UserName, &p.Title, &p.Content, &p.Media, &p.ContentType, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		p.Comments, err = GetAllCommentsForPost(p.PostId)
		if err != nil{
			return nil, err
		}
		
		posts = append(posts, p)
	}
	fmt.Println("All posts", posts)
	return posts, nil
}

func FilterCategories(categories []string) ([]Post, error) {
	posts := []Post{}
	for _, categoryID := range categories {
		var postId string
		query1 := `SELECT post_id FROM post_categories WHERE category_id = ?`
		rows, err := DB.Query(query1, categoryID)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		for rows.Next() {
			rows.Scan(&postId)
			query := `SELECT post_id, user_uuid, title, content,  media, content_type, created_at
			FROM posts p 		
			WHERE post_id = ?`

			var p Post
			err = DB.QueryRow(query, postId).Scan(&p.PostId, &p.UserId, &p.Title, &p.Content, &p.Media, &p.ContentType, &p.TimeStamp)
			if err != nil {
				return nil, err
			}
			posts = append(posts, p)
		}

		// rows, err =
		// if err != nil {
		// 	fmt.Println(err)
		// 	return nil, err
		// }
		// for rows.Next() {

		// 	err := rows.Scan(&p.Id, &p.PostId, &p.UserId, &p.Media, &p.Category, &p.Title, &p.PostContent, &p.ContentType, &p.TimeStamp)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	posts = append(posts, p)
		// }
	}
	return posts, nil
}


func MediaToBase64(media []byte)string{
	var mediaBase64 string
	if len(media) > 0 {
		mediaBase64 = base64.StdEncoding.EncodeToString(media)
	}
	return mediaBase64
}