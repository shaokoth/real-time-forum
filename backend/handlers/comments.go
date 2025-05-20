// handleComments handles CRUD operations for comments
package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"real-time-forum/backend/database"
	"real-time-forum/backend/utils"
)

func HandleComments(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	user, err := utils.GetUserFromSession(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case "GET":
		// Get comments for a post
		postID := r.URL.Query().Get("post_id")
		if postID == "" {
			http.Error(w, "Missing post ID", http.StatusBadRequest)
			return
		}

		rows, err := database.Db.Query(`
			SELECT c.id, c.content, c.user_id, c.post_id, c.created_at, u.nickname
			FROM comments c
			JOIN users u ON c.user_id = u.id
			WHERE c.post_id = ?
			ORDER BY c.created_at ASC
		`, postID)
		if err != nil {
			http.Error(w, "Error fetching comments", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Parse the comments
		var comments []map[string]interface{}
		for rows.Next() {
			var comment struct {
				ID        int
				Content   string
				UserID    int
				PostID    int
				CreatedAt time.Time
				Nickname  string
			}
			err := rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.PostID, &comment.CreatedAt, &comment.Nickname)
			if err != nil {
				http.Error(w, "Error parsing comments", http.StatusInternalServerError)
				return
			}

			// Add to the result
			comments = append(comments, map[string]interface{}{
				"id":        comment.ID,
				"content":   comment.Content,
				"userId":    comment.UserID,
				"post_id":    comment.PostID,
				"createdAt": comment.CreatedAt,
				"author":    comment.Nickname,
			})
		}

		// Return the comments
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)

	case "POST":
		// Create a new comment
		var comment struct {
			Content string `json:"content"`
			PostID  int    `json:"postId"`
		}
		err := json.NewDecoder(r.Body).Decode(&comment)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate inputs
		if comment.Content == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		// Insert the comment
		result, err := database.Db.Exec(
			"INSERT INTO comments (content, user_id, post_id) VALUES (?, ?, ?)",
			comment.Content, user.ID, comment.PostID,
		)
		if err != nil {
			http.Error(w, "Error creating comment", http.StatusInternalServerError)
			return
		}

		// Get the comment ID
		commentID, err := result.LastInsertId()
		if err != nil {
			http.Error(w, "Error getting comment ID", http.StatusInternalServerError)
			return
		}

		// Return success with the new comment ID
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":      commentID,
			"message": "Comment created successfully",
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
