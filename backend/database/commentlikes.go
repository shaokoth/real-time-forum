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
    user_id INTEGER NOT NULL,
    comment_id INTEGER NOT NULL,
    is_like BOOLEAN NOT NULL, -- true for like, false for dislike
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, comment_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (comment_id) REFERENCES comments (comment_id)
);`

	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}
