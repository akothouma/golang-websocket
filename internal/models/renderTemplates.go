package models

import (
	"html/template"
	"net/http"
)

var templates *template.Template

// Initialize templates (call this function in your main.go or init function)
func InitTemplates(templateDir string) {
	templates = template.Must(template.ParseGlob(templateDir + "/*.html"))
}

// function to render the html templates pages (used for the likes and dislike form)
// prone to change if there exists a better one
func RenderTemplates(w http.ResponseWriter, fileName string, data interface{}) {
	templates.ExecuteTemplate(w, fileName, data)
	// if err != nil {
	// 	http.Error(w, "internal server error", http.StatusInternalServerError)
	// 	log.Println("Templates failed to execute:", err)
	// 	return
	// }

	// }
}
