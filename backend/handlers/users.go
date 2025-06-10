package handlers

import (
	"encoding/json"
	"net/http"

	"real-time-forum/backend/database"
	"real-time-forum/backend/models"
	"real-time-forum/backend/utils"
)

func HandleUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	user, err := utils.GetUserFromSession(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}
	// Get all users
	rows, err := database.Db.Query(`
	SELECT id,uuid, gender, age, nickname, email, first_name, last_name
	FROM users
	ORDER BY nickname ASC
	`)
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Parse the users
	var users []models.User

	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.UUID, &u.Gender, &u.Age, &u.Nickname, &u.Email, &u.FirstName, &u.LastName)
		if err != nil {
			http.Error(w, "Error parsing users", http.StatusInternalServerError)
			return
		}
		if u.ID != user.ID {
			users = append(users, u)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
