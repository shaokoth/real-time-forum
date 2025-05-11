package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DatabasePath = "./backend/database/forum.db"
	SchemaPath   = "./backend/database/schema.sql"
)

type DB struct {
	*sql.DB
	mu sync.Mutex
}

var (
	instance *DB
	once     sync.Once
)

// Initialize sets up the database connection and schema
func Initialize() (*DB, error) {
	var err error
	once.Do(func() {
		if err = os.MkdirAll(filepath.Dir(DatabasePath), 0o755); err != nil {
			log.Printf("Error creating database directory: %v", err)
			return
		}
		var sqlDB *sql.DB
		sqlDB, err = sql.Open("sqlite3", DatabasePath+"?_foreign_keys=on")
		if err != nil {
			log.Printf("Error opening database: %v", err)
			return
		}
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(25)
		sqlDB.SetConnMaxLifetime(5 * time.Minute)
		if err = sqlDB.Ping(); err != nil {
			log.Printf("Error connecting to database: %v", err)
			return
		}
		instance = &DB{DB: sqlDB}
		if err = instance.initSchema(); err != nil {
			log.Printf("Error initializing schema: %v", err)
			return
		}
		log.Println("Database initialized successfully")
	})
	return instance, err
}

func GetInstance() (*DB, error) {
	if instance == nil {
		return Initialize()
	}
	return instance, nil
}

func (db *DB) initSchema() error {
	schema, err := os.ReadFile(SchemaPath)
	if err != nil {
		return fmt.Errorf("error reading schema file: %w", err)
	}
	_, err = db.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("error executing schema: %w", err)
	}
	return nil
}

func (db *DB) Close() error {
	if db.DB != nil {
		return db.DB.Close()
	}
	return nil
}

func (db *DB) Transaction(fn func(*sql.Tx) error) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *DB) CreateUser(nickname, email, password, firstName, lastName string, age int, gender string) (int64, error) {
	query := `INSERT INTO users (nickname, email, password, first_name, last_name, age, gender) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, nickname, email, password, firstName, lastName, age, gender)
	if err != nil {
		return 0, fmt.Errorf("error creating user: %w", err)
	}
	return result.LastInsertId()
}

func (db *DB) GetUserByEmailOrNickname(identifier string) (*sql.Row, error) {
	query := `SELECT id, nickname, email, password, first_name, last_name, age, gender 
              FROM users WHERE email = ? OR nickname = ?`
	return db.QueryRow(query, identifier, identifier), nil
}

func (db *DB) UpdateLastLogin(userID int64) error {
	query := `UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("error updating last login: %w", err)
	}
	return nil
}

func (db *DB) CreateSession(sessionID string, userID int64, expiresAt time.Time) error {
	query := `INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)`
	_, err := db.Exec(query, sessionID, userID, expiresAt)
	if err != nil {
		return fmt.Errorf("error creating session: %w", err)
	}
	return nil
}

func (db *DB) GetSessionByID(sessionID string) (*sql.Row, error) {
	query := `SELECT user_id, expires_at FROM sessions WHERE id = ?`
	return db.QueryRow(query, sessionID), nil
}

func (db *DB) DeleteSession(sessionID string) error {
	query := `DELETE FROM sessions WHERE id = ?`
	_, err := db.Exec(query, sessionID)
	if err != nil {
		return fmt.Errorf("error deleting session: %w", err)
	}
	return nil
}

func (db *DB) UpdateUserOnlineStatus(userID int64, status string) error {
	query := `INSERT INTO online_users (user_id, status, last_active) 
              VALUES (?, ?, CURRENT_TIMESTAMP) 
              ON CONFLICT(user_id) 
              DO UPDATE SET status = ?, last_active = CURRENT_TIMESTAMP`
	_, err := db.Exec(query, userID, status, status)
	if err != nil {
		return fmt.Errorf("error updating online status: %w", err)
	}
	return nil
}

func (db *DB) GetOnlineUsers() (*sql.Rows, error) {
	query := `SELECT u.id, u.nickname, u.first_name, u.last_name, o.status
              FROM users u
              JOIN online_users o ON u.id = o.user_id
              WHERE o.status = 'online' AND 
                    datetime(o.last_active) > datetime('now', '-5 minutes')
              ORDER BY u.nickname`
	return db.Query(query)
}

func (db *DB) CreatePost(userID int64, title, content string, categoryIDs []int64) (int64, error) {
	var postID int64
	err := db.Transaction(func(tx *sql.Tx) error {
		postQuery := `INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)`
		result, err := tx.Exec(postQuery, userID, title, content)
		if err != nil {
			return err
		}
		postID, err = result.LastInsertId()
		if err != nil {
			return err
		}
		for _, categoryID := range categoryIDs {
			_, err = tx.Exec(`INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`,
				postID, categoryID)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("error creating post: %w", err)
	}
	return postID, nil
}

func (db *DB) CreateComment(postID, userID int64, content string) (int64, error) {
	query := `INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)`

	result, err := db.Exec(query, postID, userID, content)
	if err != nil {
		return 0, fmt.Errorf("error creating comment: %w", err)
	}
	return result.LastInsertId()
}

func (db *DB) GetCommentsByPostID(postID int64) (*sql.Rows, error) {
	query := `SELECT c.id, c.content, c.created_at, 
              u.id as user_id, u.nickname, u.first_name, u.last_name
              FROM comments c
              JOIN users u ON c.user_id = u.id
              WHERE c.post_id = ?
              ORDER BY c.created_at ASC`
	return db.Query(query, postID)
}

func (db *DB) SendPrivateMessage(senderID, receiverID int64, content string) (int64, error) {
	query := `INSERT INTO private_messages (sender_id, receiver_id, content) VALUES (?, ?, ?)`
	result, err := db.Exec(query, senderID, receiverID, content)
	if err != nil {
		return 0, fmt.Errorf("error sending private message: %w", err)
	}
	return result.LastInsertId()
}

func (db *DB) GetConversationWithUser(userID1, userID2 int64, limit, offset int) (*sql.Rows, error) {
	query := `SELECT pm.id, pm.content, pm.created_at, pm.is_read,
              pm.sender_id, s.nickname as sender_nickname,
              pm.receiver_id, r.nickname as receiver_nickname
              FROM private_messages pm
              JOIN users s ON pm.sender_id = s.id
              JOIN users r ON pm.receiver_id = r.id
              WHERE (pm.sender_id = ? AND pm.receiver_id = ?) OR
                    (pm.sender_id = ? AND pm.receiver_id = ?)
              ORDER BY pm.created_at DESC
              LIMIT ? OFFSET ?`
	return db.Query(query, userID1, userID2, userID2, userID1, limit, offset)
}

func (db *DB) MarkMessagesAsRead(receiverID, senderID int64) error {
	query := `UPDATE private_messages 
              SET is_read = 1 
              WHERE receiver_id = ? AND sender_id = ? AND is_read = 0`
	_, err := db.Exec(query, receiverID, senderID)
	if err != nil {
		return fmt.Errorf("error marking messages as read: %w", err)
	}
	return nil
}

func (db *DB) GetUserConversations(userID int64) (*sql.Rows, error) {
	query := `WITH conversations AS (
                SELECT 
                    CASE 
                        WHEN sender_id = ? THEN receiver_id 
                        ELSE sender_id 
                    END as other_user_id,
                    MAX(created_at) as last_message_time
                FROM private_messages
                WHERE sender_id = ? OR receiver_id = ?
                GROUP BY other_user_id
              )
              SELECT 
                  u.id, u.nickname, u.first_name, u.last_name,
                  c.last_message_time,
                  (SELECT COUNT(*) FROM private_messages 
                   WHERE sender_id = u.id AND receiver_id = ? AND is_read = 0) as unread_count
              FROM conversations c
              JOIN users u ON c.other_user_id = u.id
              LEFT JOIN online_users o ON u.id = o.user_id
              ORDER BY c.last_message_time DESC, u.nickname ASC`
	return db.Query(query, userID, userID, userID, userID)
}