package database

import (
	"database/sql"
	"fmt"
)

// Private messages Table
func CreatePrivateMessagesTable(db *sql.DB) error {
	if db == nil {
		// 	defer db.Close()
		return fmt.Errorf("nil database connection")
	}
	Query := `
	CREATE TABLE IF NOT EXISTS private_messages (
	private_id INTEGER PRIMARY KEY AUTOINCREMENT,
	content TEXT NOT NULL,
    sender_id INTEGER NOT NULL,
    receiver_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_read BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (sender_id) REFERENCES users(id),
    FOREIGN KEY (receiver_id) REFERENCES users(id)
	)`
	if _, err := db.Exec(Query); err != nil {
		return err
	}
	return nil
}
