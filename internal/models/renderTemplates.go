package models

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

var templates *template.Template

// Initialize templates (call this function in your main.go or init function)
func InitTemplates() (*template.Template, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, errors.New("error getting working directory")
	}
	// constructing the correct path
	tmplPath := wd + "/ui/templates/"
	tmpl, err := template.ParseFiles(
		tmplPath+"base.html",
		tmplPath+"categories.html",
		tmplPath+"home.html",
		tmplPath+"postContent.html",
		tmplPath+"posts.html",
	)
	if err != nil {
		return nil, fmt.Errorf("error parsing templates: %w", err) // Wrap the error
	}
	return tmpl, nil
}

// function to render the html templates pages (used for the likes and dislike form)
// prone to change if there exists a better one
func RenderTemplates(w http.ResponseWriter, fileName string, data interface{}) {
	if err := templates.ExecuteTemplate(w, fileName, data); err != nil {
		ErrorHandler(w, http.StatusInternalServerError)
		log.Println("Templates failed to execute:", err)
		return
	}
}
