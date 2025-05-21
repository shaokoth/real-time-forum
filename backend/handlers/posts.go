package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"real-time-forum/backend/database"
	"real-time-forum/backend/utils"
)

// Handles CRUD operations for posts
func HandlePosts(w http.ResponseWriter, r *http.Request) {
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
	switch r.Method {
	case "GET":
		// Get all posts
		rows, err := database.Db.Query(`
			SELECT p.id, p.title, p.content, p.user_id, p.category_id, p.created_at, u.nickname, c.name
			FROM posts p
			JOIN users u ON p.user_id = u.id
			JOIN categories c ON p.category_id = c.id
			ORDER BY p.created_at DESC
		`)
		if err != nil {
			http.Error(w, "Error fetching posts", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Parse the posts
		var posts []map[string]interface{}
		for rows.Next() {
			var post struct {
				ID           int
				Title        string
				Content      string
				UserID       int
				CategoryID   int
				CreatedAt    time.Time
				Nickname     string
				CategoryName string
			}
			err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CategoryID, &post.CreatedAt, &post.Nickname, &post.CategoryName)
			if err != nil {
				http.Error(w, "Error parsing posts", http.StatusInternalServerError)
				return
			}

			// Add to the result
			posts = append(posts, map[string]interface{}{
				"id":           post.ID,
				"title":        post.Title,
				"content":      post.Content,
				"userId":       post.UserID,
				"categoryId":   post.CategoryID,
				"createdAt":    post.CreatedAt,
				"author":       post.Nickname,
				"categoryName": post.CategoryName,
			})
		}

		// Return the posts
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(posts)
		
	case "POST":
		// Create a new post
		var post struct {
			Title      string `json:"title"`
			Content    string `json:"content"`
			Categories []int  `json:"categories"`
		}
		err := json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Validate inputs
		if post.Title == "" || post.Content == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		if len(post.Categories) == 0 {
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
			"INSERT INTO posts (title, content, user_id) VALUES (?, ?, ?)",
			post.Title, post.Content, user.ID,
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
		for _, categoryID := range post.Categories {
			_, err := tx.Exec(
				"INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)",
				postID, categoryID,
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
