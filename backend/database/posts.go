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
        post_id INTEGER  PRIMARY KEY AUTOINCREMENT DEFAULT 0,
		user_uuid TEXT NOT NULL,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        filepath TEXT DEFAULT '',
        filename TEXT DEFAULT '',
        category  TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_uuid) REFERENCES users(uuid)
    );`

	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}
