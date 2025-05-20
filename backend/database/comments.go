package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// ==== The function will create the comments table in the database =====
func CreateCommentsTable(db *sql.DB) error {
	if db == nil {
		// defer db.Close()
		// e.LOGGER("[ERROR]", fmt.Errorf("|create comments table|"))
		return fmt.Errorf("nil database connection")
	}

	query := `
    CREATE TABLE IF NOT EXISTS comments (
        comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_uuid TEXT NOT NULL,
        post_id INTEGER NOT NULL,
        content TEXT NOT NULL,
        likes INTEGER DEFAULT 0,
        dislikes INTEGER DEFAULT 0,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_uuid) REFERENCES posts(uuid)
        FOREIGN KEY (post_id) REFERENCES posts(post_id)

    );`

	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}
