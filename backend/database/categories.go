package database

import (
	"database/sql"
	"fmt"
)

// ==== The function will create the categories table in the database =====
func CreateCategoriesTable(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("nil database connection")
	}

	query := `
    CREATE TABLE IF NOT EXISTS categories (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL
    );`

	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}
