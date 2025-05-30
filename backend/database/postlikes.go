package database

import (
	"database/sql"
	"fmt"
)

func CreatePostLikes(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("nil database connection")
	}
	query := `CREATE TABLE IF NOT EXISTS post_likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    is_like BOOLEAN NOT NULL, -- true for like, false for dislike
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, post_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (post_id) REFERENCES posts (id)
);`

	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}
