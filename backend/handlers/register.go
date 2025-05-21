package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"real-time-forum/backend/database"
	"real-time-forum/backend/models"
	"real-time-forum/backend/utils"

	"github.com/gofrs/uuid"
)

// Handles user registration
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}
	 
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if !utils.IsValidEmail(user.Email) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid email address"))
		//http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	if utils.CredentialExists(database.Db, user.Nickname) || utils.CredentialExists(database.Db, user.Email) {
		//w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("nickname or email already exists"))
		//http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	if err := user.HashPassword(); err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	u, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	UUID := u.String()

	_, err = database.Db.Exec(`INSERT INTO users(uuid, nickname,age,gender,first_name,email,last_name,password)VALUES(?,?,?,?,?,?,?,?)`, UUID, user.Nickname, user.Age, user.Gender, user.FirstName, user.Email, user.LastName, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var userID string
	err = database.Db.QueryRow("SELECT id FROM users WHERE email = ?", user.Email).Scan(&userID)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}
