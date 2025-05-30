package database

import (
	"database/sql"
	"fmt"
)

func CreateCommentLikes(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("nil database connection")
	}
	query := `CREATE TABLE IF NOT EXISTS comment_likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    comment_id INTEGER NOT NULL,
    is_like BOOLEAN NOT NULL, -- true for like, false for dislike
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, comment_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (comment_id) REFERENCES comments (id)
);`

	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}
