package database

import (
	"database/sql"
	"fmt"
)

// ==== The function will create the categories table in the database =====
func CreatePostCategoriesTable(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("nil database connection")
	}

	query := `
    CREATE TABLE IF NOT EXISTS post_categories (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        post_id INTEGER NOT NULL,
        UNIQUE(post_id, name),
        FOREIGN KEY(post_id) REFERENCES posts(post_id)
    );`

	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}
