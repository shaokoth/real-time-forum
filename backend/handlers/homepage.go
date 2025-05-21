package handlers

import (
	"html/template"
	"net/http"
)

// HandleHomepage serves the index.html page for the forum platform
func HandleHomepage(w http.ResponseWriter, r *http.Request) {
	// Only handle the root path
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Get user session info
	
	// Parse and serve the template
	tmpl, err := template.ParseFiles("/home/docker/real-time-forum/frontend/template/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
