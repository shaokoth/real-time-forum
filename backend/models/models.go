package models

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

// Client represents a connected client
type Client struct {
	UserID     int
	Connection *websocket.Conn
	Send       chan []byte
}

type User struct {
	ID        int       `json:"id"`
	UUID      string    `json:"user_uuid"`
	Nickname  string    `json:"nickname"`
	Age       int       `json:"age"`
	Gender    string    `json:"gender"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"Password"`
	CreatedAt time.Time `json:"createdAt"`
	IsOnline  bool      `json:"isOnline"`
}

type Post struct {
	CreatedAt     time.Time `json:"created_at"`
	Categories    []string  `json:"category"`
	Likes         int       `json:"likes"`
	Title         string    `json:"title"`
	Dislikes      int       `json:"dislikes"`
	CommentsCount int       `json:"comments_count"`
	Comments      []Comment `json:"comments"`
	Content       string    `json:"content"`
	User_uuid     string    `json:"user_uuid"`
	Post_id       int       `json:"post_id"`
	Filepath      string    `json:"filepath"`
	Filename      string    `json:"filename"`
	Owner         string
	OwnerInitials string
}

type Comment struct {
	Comment_id int       `json:"comment_id"`
	Post_id    string    `json:"post_id"`
	CreatedAt  time.Time `json:"created_at"`
	Likes      int       `json:"likes"`
	Dislikes   int       `json:"dislikes"`
	Content    string    `json:"content"`
}

type PrivateMessage struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// =====  hashes the user's password before storing it ====
func (user *User) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("|hashpassword method| ---> {%v}", err)
		return err
	}
	user.Password = string(hashed)
	return nil
}
