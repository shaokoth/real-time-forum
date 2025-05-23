package handlers

import (
	"encoding/json"
	"net/http"
	"sync"

	"real-time-forum/backend/database"
	"real-time-forum/backend/models"
	"real-time-forum/backend/utils"
)

var mu sync.Mutex
var activeUsers map[string]*models.User

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
	var users []map[string]interface{}
	for rows.Next() {
		var user struct {
			ID        int
			Nickname  string
			Email     string
			FirstName string
			LastName  string
		}
		err := rows.Scan(&user.ID, &user.Nickname, &user.Email, &user.FirstName, &user.LastName)
		if err != nil {
			http.Error(w, "Error parsing users", http.StatusInternalServerError)
		}

		// Check if user is online
		mu.Lock()
		_, isOnline := activeUsers[user.Nickname]
		mu.Unlock()
		users = append(users, map[string]interface{}{
			"id":        user.ID,
			"nickname":  user.Nickname,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"isOnline":  isOnline,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
