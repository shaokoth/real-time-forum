package handlers

import (
	"encoding/json"
	"net/http"
	"sync"

	"real-time-forum/backend/database"
	"real-time-forum/backend/models"
	"real-time-forum/backend/utils"
)

var (
	mu          sync.Mutex
	users       []models.User
	// activeUsers map[string]*models.User
)

func HandleUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	_, err = utils.GetUserFromSession(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}
	// Get all users
	rows, err := database.Db.Query(`
	SELECT id, nickname, email, firstname, lastname
	FROM users
	ORDER BY nickname ASC
	`)
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Parse the users
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Nickname, &user.Email, &user.FirstName, &user.LastName)
		if err != nil {
			http.Error(w, "Error parsing users", http.StatusInternalServerError)
		}

	users = append(users, user)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
