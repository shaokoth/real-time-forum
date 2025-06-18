package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"real-time-forum/backend/database"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var DefaultCategories = []Category{
	{ID: 1, Name: "Sports"},
	{ID: 2, Name: "Lifestyle"},
	{ID: 3, Name: "Education"},
	{ID: 4, Name: "Finance"},
	{ID: 5, Name: "Music"},
	{ID: 6, Name: "Culture"},
	{ID: 7, Name: "ReactedPosts"},
}

// GetAllCategories returns all categories from the database or defaults if none exist
func GetAllCategories() ([]Category, error) {
	query := `SELECT id, name FROM Categories ORDER BY name`
	rows, err := database.Db.Query(query)
	if err != nil {
		return DefaultCategories, nil
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return nil, fmt.Errorf("failed to scan category: %v", err)
		}
		categories = append(categories, cat)
	}

	// If no categories in database, return defaults
	if len(categories) == 0 {
		return DefaultCategories, nil
	}

	return categories, nil
}

// HandleGetCategories handles GET requests for retrieving all categories
func HandleGetCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	categories, err := GetAllCategories()
	if err != nil {
		http.Error(w, "Error fetching categories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
