package main

import (
	"net/http"

	"learn.zone01kisumu.ke/git/clomollo/forum/internal/models"
	// middleware "learn.zone01kisumu.ke/git/clomollo/forum/Middleware"
)

func (dep *Dependencies) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	models.InitTemplates("./ui/html/")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// mux.Handle("/", dep.CSRFMiddleware(http.HandlerFunc(dep.HomeHandler)))
	mux.Handle("/post", dep.AuthMiddleware(http.HandlerFunc(dep.PostHandler)))
	// mux.Handle("/allposts", http.HandlerFunc(dep.AllPostsHandler))
	mux.Handle("/allposts", http.HandlerFunc(models.RenderPostsPage))
	
	

	mux.Handle("/register", dep.CSRFMiddleware(http.HandlerFunc(dep.RegisterHandler)))
	mux.Handle("/logout", http.HandlerFunc(dep.LogoutHandler))
	mux.Handle("/login", dep.CSRFMiddleware(http.HandlerFunc(dep.LoginHandler)))
	mux.Handle("/styling/", http.StripPrefix("/styling/", http.FileServer(http.Dir("./ui/static/styling"))))
	mux.Handle("/add_comment", dep.AuthMiddleware(http.HandlerFunc(dep.AddCommentHandler)))
	mux.Handle("/filtered_posts",http.HandlerFunc(dep.PostsByFilters))

	mux.Handle("/likes", dep.AuthMiddleware(http.HandlerFunc(dep.LikeHandler)))

	return mux
}
