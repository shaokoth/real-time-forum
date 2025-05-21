package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// ==== The function will create the posts table in the database =====
func CreatePostsTable(db *sql.DB) error {
	if db == nil {
		// 	defer db.Close()
		return fmt.Errorf("nil database connection")
	}

	query := `
	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}
