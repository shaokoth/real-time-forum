package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"real-time-forum/backend/database"
	"real-time-forum/backend/models"
	"real-time-forum/backend/utils"
)

// Handles CRUD operations for posts
func HandlePosts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// Get all posts - no authentication required
		var posts []models.Post
		rows, err := database.Db.Query(`
			SELECT p.post_id, p.title, p.content, p.user_uuid, p.created_at,
				   GROUP_CONCAT(pc.name) as categories,
				   COALESCE(SUM(CASE WHEN pl.is_like = 1 THEN 1 ELSE 0 END), 0) as likes,
				   COALESCE(SUM(CASE WHEN pl.is_like = 0 THEN 1 ELSE 0 END), 0) as dislikes
			FROM posts p
			LEFT JOIN post_categories pc ON p.post_id = pc.post_id
			LEFT JOIN post_likes pl ON p.post_id = pl.post_id
			GROUP BY p.post_id
			ORDER BY p.created_at DESC
		`)
		if err != nil {
			http.Error(w, "Error fetching posts", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Parse the posts
		for rows.Next() {
			var post models.Post
			var categoriesStr string
			err = rows.Scan(&post.Post_id, &post.Title, &post.Content, &post.User_uuid, &post.CreatedAt, &categoriesStr, &post.Likes, &post.Dislikes)
			if err != nil {
				http.Error(w, "Error parsing posts", http.StatusInternalServerError)
				return
			}

			// Split categories string into array
			if categoriesStr != "" {
				post.Categories = strings.Split(categoriesStr, ",")
			} else {
				post.Categories = []string{}
			}

			posts = append(posts, post)
		}
		// Return the posts
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)

	case "POST":
		// Create a new post - authentication required
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

		var post struct {
			Title   string `json:"title"`
			Content string `json:"content"`
			Numbers []int  `json:"categories"`
		}
		err = json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		var Categories []Category
		for _, v := range post.Numbers {
			Categories = append(Categories, DefaultCategories[v-1])
		}
		// Validate inputs
		if post.Title == "" || post.Content == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		if len(Categories) == 0 {
			http.Error(w, "At least one category is required", http.StatusBadRequest)
			return
		}

		// Start a transaction
		tx, err := database.Db.Begin()
		if err != nil {
			http.Error(w, "Error starting transaction", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		// Insert the post
		result, err := tx.Exec(
			"INSERT INTO posts (title, content, user_uuid) VALUES (?, ?, ?)",
			post.Title, post.Content, user.UUID,
		)
		if err != nil {
			http.Error(w, "Error creating post", http.StatusInternalServerError)
			return
		}

		// Get the post ID
		postID, err := result.LastInsertId()
		if err != nil {
			http.Error(w, "Error getting post ID", http.StatusInternalServerError)
			return
		}
		
		// Insert post categories
		for _, name := range Categories {
			_, err := tx.Exec(
				"INSERT INTO post_categories (post_id, name) VALUES ( ?, ?)",
				postID, name.Name,
			)
			if err != nil {
				http.Error(w, "Error adding category to post", http.StatusInternalServerError)
				return
			}
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			http.Error(w, "Error committing transaction", http.StatusInternalServerError)
			return
		}

		// Return success with the new post ID
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":      postID,
			"message": "Post created successfully",
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	type LikeRequest struct {
		UserID int  `json:"user_id"`
		PostID int  `json:"post_id"`
		IsLike bool `json:"is_like"` // true = like, false = dislike
	}

	// Get user ID from session
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
	// Parse request body
	var req LikeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	// Check if user already reacted to this post
	existingReaction, err := utils.GetPostReaction(user.ID, req.PostID)
	if err == nil {
		if existingReaction.IsLike {
			// User already liked - remove like
			err = utils.DeletePostReaction(user.ID, req.PostID)
			if err != nil {
				http.Error(w, "Failed to remove like", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Like removed"))
			return
		} else {
			// User disliked - change to like
			err = utils.UpdatePostReaction(user.ID, req.PostID, req.IsLike)
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
	err = utils.AddPostReaction(user.ID, req.PostID, req.IsLike)
	if err != nil {
		http.Error(w, "Failed to add like", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Post liked"))
}

func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	type LikeRequest struct {
		UserID int  `json:"user_id"`
		PostID int  `json:"post_id"`
		IsLike bool `json:"is_like"` // true = like, false = dislike
	}

	// Get user ID from session
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
	// Parse request body
	var req LikeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	// Check if user already reacted to this post
	existingReaction, err := utils.GetPostReaction(user.ID, req.PostID)
	if err == nil {
		if existingReaction.IsLike {
			// User already disliked - remove dislike
			err = utils.DeletePostReaction(user.ID, req.PostID)
			if err != nil {
				http.Error(w, "Failed to remove dislike", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Dislike removed"))
			return
		} else {
			// User disliked - change to like
			err = utils.UpdatePostReaction(user.ID, req.PostID, false)
			if err != nil {
				http.Error(w, "Failed to update reaction", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("like changed to dislike"))
			return
		}
	}

	// Add new dislike
	err = utils.AddPostReaction(user.ID, req.PostID, false)
	if err != nil {
		http.Error(w, "Failed to add dislike", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Post disliked"))
}
