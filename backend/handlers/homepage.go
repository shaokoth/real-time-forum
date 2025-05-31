package handlers

import (
	"html/template"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/utils"
)

// HandleHomepage serves the index.html page for the forum platform
func HandleHomepage(w http.ResponseWriter, r *http.Request) {
	// Only handle the root path
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Check if user is logged in
	isLoggedIn := false
	var nickname string
	valid, userID := utils.ValidateSession(r)
	if valid {
		isLoggedIn = true
		// Get user nickname
		err := database.Db.QueryRow("SELECT nickname FROM users WHERE id = ?", userID).Scan(&nickname)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	// Parse and serve the template
	tmpl, err := template.ParseFiles("frontend/template/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create template data
	data := struct {
		IsLoggedIn bool
		Nickname   string
	}{
		IsLoggedIn: isLoggedIn,
		Nickname:   nickname,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
