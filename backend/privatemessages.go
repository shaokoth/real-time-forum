package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// ==== The function will create the privatemessages table in the database =====
func  CreatePrivateMessages(db *sql.DB) error {
	if db == nil {
		// defer db.Close()
		// e.LOGGER("[ERROR]", fmt.Errorf("|create comments table|"))
		return fmt.Errorf("nil database connection")
	}

	query := `
    CREATE TABLE IF NOT EXISTS privatemessages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_uuid TEXT NOT NULL,
    sender_id INTEGER NOT NULL,
    receiver_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_read BOOLEAN DEFAULT 0,
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE

    );`

	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}
