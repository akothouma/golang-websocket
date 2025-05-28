package main

import (
	"log"
	"net/http"
	"os"

	handlers "learn.zone01kisumu.ke/git/clomollo/forum/internal/Handlers"
	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
)

func Routes() *http.ServeMux {
	// leveled logging, informational messages output to standard out(stdout)
	// Error messages output to standard error(stderr)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t",
		log.Ldate|log.Ltime|log.Lshortfile)

	db, err := models.InitializeDB()
	if err != nil {
		errorLog.Fatalf("Failed to initialize database: %v", err)
	}

	models.DB = db

	// initializing dependencies
	dep := &handlers.Dependencies{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
		Forum:    &models.ForumModel{DB: db},
	}
    go dep.BroadcastToClients();
	mux := http.NewServeMux()

	models.InitTemplates("./ui/html/")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static"))))
	mux.Handle("/post", dep.AuthMiddleware(http.HandlerFunc(dep.PostHandler)))
	mux.Handle("/", http.HandlerFunc(models.RenderPostsPage))
	mux.HandleFunc("/my_posts", models.RenderMyPostsPage)
	mux.HandleFunc("/liked_posts", models.RenderLikedPostsPage)
	mux.Handle("/upload-profile", dep.AuthMiddleware(http.HandlerFunc(handlers.UploadProfilePictureHandler)))
	mux.Handle("/register", dep.CSRFMiddleware(http.HandlerFunc(dep.RegisterHandler)))
	mux.Handle("/ws",dep.AuthMiddleware(http.HandlerFunc(dep.ChatHandler)))
	mux.Handle("/logout", http.HandlerFunc(dep.LogoutHandler))
	mux.Handle("/login", dep.CSRFMiddleware(http.HandlerFunc(dep.LoginHandler)))
	mux.Handle("/add_comment", dep.AuthMiddleware(http.HandlerFunc(handlers.AddCommentHandler)))
	mux.Handle("/add_reply", dep.AuthMiddleware(http.HandlerFunc(handlers.AddReplyHandler)))
	mux.Handle("/filtered_posts", http.HandlerFunc(handlers.PostsByFilters))
	mux.Handle("/likes", dep.AuthMiddleware(http.HandlerFunc(dep.LikeHandler)))

	return mux
}
