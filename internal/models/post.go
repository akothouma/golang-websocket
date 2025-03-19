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
	PostId         string         `json:"PostId"`
	UserId         string         `json:"UserId"`
	Media          []byte         `json:"Media"`
	MediaString    string         `json:"MediaString"`
	Category       []string       `json:"Category"`
	Title          string         `json:"Title"`
	Content        string         `json:"Content"`
	ContentType    string         `json:"ContentType"`
	Likes          int            `json:"Likes"`
	Dislikes       int            `json:"Dislikes"`
	Comments       []Comment      `json:"Comments"`
	CommentsLenght int            `json:"CommentsLenght"`
	UserName       string         `json:"UserName"`
	Initial        string         `json:"Initial"`
	Categories     []postCategory `json:"Categories"`
	CreatedAt      time.Time      `json:"CreatedAt"`
}

// type LikeData struct {
//     Count    int
//     Reaction bool
// }

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
            INSERT INTO post_categories (postId, category_id)
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

		p.Initial = string(p.UserName[0])

		p.MediaString = MediaToBase64(p.Media)

		p.Comments, err = GetAllCommentsForPost(p.PostId)
		if err != nil {
			return nil, err
		}

		p.Likes, p.Dislikes, err = PostLikesDislikes(p.PostId)
		if err != nil {
			return nil, err
		}

		p.Categories, err = Post_Categories(p.PostId)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}
	return posts, nil
}

func FilterCategories(categories []string) ([]string, error) {
	// posts := []Post{}
	data := []string{}

	for _, categoryID := range categories {
		id := ""

		query := `SELECT p.post_id
		FROM posts p 
		JOIN post_categories pc ON pc.post_id = p.post_id		
		WHERE pc.category_id = ?`

		err := DB.QueryRow(query, categoryID).Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows{
				// fmt.Println("FilterCategories Err:", err)
				continue
			}else{
				return nil, err
			}			
		}

		data = append(data, id)


		// query := `SELECT p.id, p.post_id, p.user_uuid, p.username, p.title, p.content, p.media, p.content_type, p.created_at 
		// FROM posts p 
		// JOIN post_categories pc ON pc.post_id = p.post_id		
		// WHERE pc.category_id = ?`

		// var p Post
		// err := DB.QueryRow(query, categoryID).Scan(&p.ID, &p.PostId, &p.UserId, &p.UserName, &p.Title, &p.Content, &p.Media, &p.ContentType, &p.CreatedAt)
		// if err != nil {
		// 	fmt.Println("FilterCategories Err:", err)
		// 	return nil, err
		// }
		// p.Initial = string(p.UserName[0])

		// p.MediaString = MediaToBase64(p.Media)

		// p.Comments, err = GetAllCommentsForPost(p.PostId)
		// if err != nil {
		// 	return nil, err
		// }

		// p.Likes, p.Dislikes, err = PostLikesDislikes(p.PostId)
		// if err != nil {
		// 	return nil, err
		// }

		// p.Categories, err = Post_Categories(p.PostId)
		// if err != nil {
		// 	return nil, err
		// }

		// posts = append(posts, p)

	}
	// return posts, nil
	return data, nil
}

func MediaToBase64(media []byte) string {
	var mediaBase64 string
	if len(media) > 0 {
		mediaBase64 = base64.StdEncoding.EncodeToString(media)
	}
	return mediaBase64
}
