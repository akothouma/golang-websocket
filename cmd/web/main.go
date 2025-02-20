package main

import (
	"flag"
	"log"
	"net/http"

	middleware "learn.zone01kisumu.ke/git/clomollo/forum/Middleware"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/database"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/handlers"
)

func main() {
	addr := flag.String("addr", ":8000", "HTTP network address")
	if err := database.InitializeDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.Handle("/register", middleware.CSRFMiddleware(http.HandlerFunc(handlers.RegisterHandler)))
	mux.Handle("/logout", http.HandlerFunc(handlers.LogoutHandler))
	mux.Handle("/login", middleware.CSRFMiddleware(http.HandlerFunc(handlers.LoginHandler)))
	mux.Handle("/add_comment", middleware.CSRFMiddleware(http.HandlerFunc(handlers.AddCommentHandler)))
	mux.Handle("/get_all_post_comment", middleware.CSRFMiddleware(http.HandlerFunc(handlers.GetAllCommentsForPostHandler)))
	mux.Handle("/get_all_comment_replies", middleware.CSRFMiddleware(http.HandlerFunc(handlers.GetAllRepliesForCommentHandler)))
	mux.Handle("/add__post", http.HandlerFunc(handlers.AddPostHandler))
	mux.Handle("/get_all_posts", http.HandlerFunc(handlers.GetAllPostsHandler))

	serv := &http.Server{
		Handler: mux,
		Addr:    *addr,
	}

	log.Printf("Starting server on port %v:", *addr)
	err := serv.ListenAndServe()
	log.Fatal(err)
}
