package models

import (
	"database/sql"
	"fmt"
	"log"
)

type Post struct {
	Id          int
	PostId      string
	UserId      string
	Media       []byte
	Category    []string
	Title       string
	PostContent string
	ContentType string
	TimeStamp   string
}

func (f *ForumModel) CreatePost(p *Post) error {
	var username string
	err := DB.QueryRow("SELECT username FROM users WHERE user_uuid = ?", p.UserId).Scan(&username)
	if err != nil {
		log.Printf("Failed to fetch username: %v", err)
		return fmt.Errorf("failed to fetch username: %w", err)
	}

	query := "INSERT INTO posts(post_id, user_uuid, username, title, content,  media, content_type) VALUES(?, ?, ?, ?, ?, ?, ?)"
	_, err = f.DB.Exec(query, p.PostId, p.UserId, username, p.Title, p.PostContent, p.Media, p.ContentType)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("failed to insert a post")
	}

	// Insert categories
	fmt.Println("categories",p.Category)
	
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

func (f *ForumModel) FindPostById(id string) (*Post, error) {
	query := "SELECT id,user_uuid,category,title,content FROM posts WHERE id=?"
	row := f.DB.QueryRow(query, id)
	post := Post{}
	err := row.Scan(&post.PostId, &post.UserId, &post.Title, &post.PostContent)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &post, nil
}

func (f *ForumModel) AllPosts() ([]Post, error) {
	query := "SELECT * FROM posts"
	rows, err := f.DB.Query(query)
	posts := []Post{}
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.Id, &p.PostId, &p.UserId, &p.Title, &p.PostContent, &p.TimeStamp)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	fmt.Println("All posts",posts)
	return posts, nil
}

func (f *ForumModel) FilterCategories(categories []string) ([]Post, error) {
	posts := []Post{}
	for _, categoryID := range categories {
		var postId string
		query1:=`SELECT post_id FROM post_categories WHERE category_id = ?`
		rows, err := f.DB.Query(query1, categoryID)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		for rows.Next(){
			rows.Scan( &postId)
			query := `SELECT post_id, user_uuid, title, content,  media, content_type, created_at
			FROM posts p 		
			WHERE post_id = ?`

			var p Post
			err = f.DB.QueryRow(query, postId).Scan(&p.PostId, &p.UserId,  &p.Title,  &p.PostContent, &p.Media, &p.ContentType, &p.TimeStamp)
			
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