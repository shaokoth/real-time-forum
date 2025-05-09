package database

import (
	"database/sql"
	"fmt"
)

var Db *sql.DB

func InitDB() error {
	Db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	err = Db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	if err = CreateUsersTable(Db); err != nil {
		return err
	}

	if err = CreateLikesDislikesTable(Db); err != nil {
		return err
	}

	if err = CreateSessionsTable(Db); err != nil {
		return err
	}

	if err = CreatePostsTable(Db); err != nil {
		return err
	}
	if err = CreateCommentsTable(Db); err != nil {
		return err
	}

	fmt.Println("SUCCESS: Connected to the SQLite database!")
	return nil
}
