package models

import (
	"database/sql"
	"fmt"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/database"
)

type Post struct {
	PostId      string   `json:"id"`
	UserId      string   `json:"user_id"`
	Category    []string `json:"category"`
	Title       string   `json:"title"`
	PostContent string   `json:"content_type"`
	Media       string   `json:"videoLink"`
}

func CreatePost(id, title, postContent string, category []string) error {
	query := "INSERT INTO posts(id,category,title,content) VALUES(?,?,?,?)"
	_, err := database.DB.Exec(query, id, title, postContent, category)
	if err != nil {
		return fmt.Errorf("failed to insert a post")
	}
	return nil
}
func FindPostById(id string)(*Post,error){
query:="SELECT id,category,title,content FROM posts WHERE id=?"
row:=database.DB.QueryRow(query,id)
post:=Post{}
err:=row.Scan(&post.PostId,&post.Category,&post.Title,&post.PostContent)
if err !=nil{
	if err==sql.ErrNoRows{
		return nil,err
	}
	return nil,err
}
return &post,nil
}