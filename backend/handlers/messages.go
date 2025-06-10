package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"real-time-forum/backend/database"
	"real-time-forum/backend/models"
	"real-time-forum/backend/utils"
)

func HandleGetMessages(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	otherUserID := r.URL.Query().Get("with")
	if otherUserID == "" {
		http.Error(w, "Missing 'with' parameter", http.StatusBadRequest)
		return
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}
	// Get current user from session
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

	// Query messages
	rows, err := database.Db.Query(`
		SELECT sender_id, receiver_id, content, created_at
		FROM private_messages
		WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
		ORDER BY created_at ASC
		LIMIT 10`,
		otherUserID, user.UUID, user.UUID, otherUserID, offset)
		
	if err != nil {
		log.Printf("Message query error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.CreatedAt); err != nil {
			log.Printf("Message scan error: %v", err)
			continue
		}
		messages = append(messages, msg)
	}

	if messages == nil {
		messages = []models.Message{}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
