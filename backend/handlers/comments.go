// handleComments handles CRUD operations for comments
package handlers

import (
	"encoding/json"
	"net/http"

	"real-time-forum/backend/database"
	"real-time-forum/backend/models"
	"real-time-forum/backend/utils"
)

func HandleComments(w http.ResponseWriter, r *http.Request) {
	var comments []models.Comment
	var comment models.Comment
	switch r.Method {
	case "GET":
		// Get comments for a post
		postID := r.URL.Query().Get("post_id")
		if postID == "" {
			http.Error(w, "Missing post ID", http.StatusBadRequest)
			return
		}

		// Get user ID if logged in
		var userID int
		cookie, err := r.Cookie("session_token")
		if err == nil {
			user, err := utils.GetUserFromSession(cookie.Value)
			if err == nil {
				userID = user.ID
			}
		}

		rows, err := database.Db.Query(`
			SELECT 
				c.comment_id, c.user_uuid, c.post_id, c.content, c.created_at, u.nickname as author,
				COALESCE(SUM(CASE WHEN cl.is_like = 1 THEN 1 ELSE 0 END), 0) as likes,
				COALESCE(SUM(CASE WHEN cl.is_like = 0 THEN 1 ELSE 0 END), 0) as dislikes,
				CASE 
					WHEN EXISTS (
						SELECT 1 FROM comment_likes WHERE comment_id = c.comment_id AND user_id = ? AND is_like = 1
					) THEN 1
					WHEN EXISTS (
						SELECT 1 FROM comment_likes 
						WHERE comment_id = c.comment_id AND user_id = ? AND is_like = 0
					) THEN -1
					ELSE 0
				END as user_liked
			FROM comments c
			JOIN users u ON c.user_uuid = u.uuid
			LEFT JOIN comment_likes cl ON c.comment_id = cl.comment_id
			WHERE c.post_id = ?
			GROUP BY c.comment_id  ORDER BY c.created_at ASC
		`, userID, userID, postID)
		if err != nil {
			http.Error(w, "Error fetching comments", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Parse the comments
		for rows.Next() {
			var comment models.Comment
			err := rows.Scan(
				&comment.Comment_id, &comment.User_uuid, &comment.Post_id, &comment.Content, &comment.CreatedAt,
				&comment.Author, &comment.Likes, &comment.Dislikes, &comment.UserLiked,
			)
			if err != nil {
				http.Error(w, "Error parsing comments", http.StatusInternalServerError)
				return
			}
			comments = append(comments, comment)
		}
		// Return the comments
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)

	case "POST":
		// Check if user is logged in
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
		err = json.NewDecoder(r.Body).Decode(&comment)
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
			"INSERT INTO comments (user_uuid, post_id, content) VALUES (?, ?, ?)",
			user.UUID, comment.Post_id, comment.Content,
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

// This function will handle liking a comment
func LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		CommentID int `json:"comment_id"`
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	user, err := utils.GetUserFromSession(cookie.Value)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	// Check if user already reacted to this comment
	existingReaction, err := utils.GetCommentReaction(user.ID, req.CommentID)
	if err == nil {
		if existingReaction.IsLike {
			// User already liked - remove like
			err = utils.DeleteCommentReaction(user.ID, req.CommentID)
			if err != nil {
				http.Error(w, "Failed to remove like", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Like removed"))
			return
		} else {
			// User disliked - change to like
			err = utils.UpdateCommentReaction(user.ID, req.CommentID, true)
			if err != nil {
				http.Error(w, "Failed to update reaction", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Dislike changed to like"))
			return
		}
	}

	// Add new like
	err = utils.AddCommentReaction(user.ID, req.CommentID, true)
	if err != nil {
		http.Error(w, "Failed to add like", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Comment liked"))
}

func DislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		CommentID int `json:"comment_id"`
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	user, err := utils.GetUserFromSession(cookie.Value)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	// Check if user already reacted to this comment
	existingReaction, err := utils.GetCommentReaction(user.ID, req.CommentID)
	if err == nil {
		if existingReaction.IsLike {
			// User disliked - change to like
			err = utils.UpdateCommentReaction(user.ID, req.CommentID, false)
			if err != nil {
				http.Error(w, "Failed to update reaction", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Like changed to dislike"))
			return
		} else {
			// User already disliked - remove dislike
			err = utils.DeleteCommentReaction(user.ID, req.CommentID)
			if err != nil {
				http.Error(w, "Failed to remove dislike", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Dislike removed"))
			return
		}
	}
	// Add new like
	err = utils.AddCommentReaction(user.ID, req.CommentID, false)
	if err != nil {
		http.Error(w, "Failed to add dislike", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Comment disliked"))
}
