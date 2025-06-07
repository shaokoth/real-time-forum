package models

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Message struct {
	ID         string    `json:"id"`
	Content    string    `json:"content"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type Reaction struct {
	ID        int  `json:"id"`
	UserID    int  `json:"user_id"`
	PostID    int  `json:"post_id,omitempty"`
	CommentID int  `json:"comment_id,omitempty"`
	IsLike    bool `json:"is_like"`
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

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Post struct {
	CreatedAt     time.Time `json:"created_at"`
	Categories    []string  `json:"categories"`
	Category      string    `json:"Category"`
	Likes         int       `json:"likes"`
	Title         string    `json:"title"`
	Dislikes      int       `json:"dislikes"`
	UserLiked     int       `json:"UserLiked"` // -1: dislike, 0: none, 1:like
	CommentsCount int       `json:"comments_count"`
	Comments      []Comment `json:"comments"`
	Content       string    `json:"content"`
	User_uuid     string    `json:"user_uuid"`
	Post_id       int       `json:"post_id"`
	Filepath      string    `json:"filepath"`
	Filename      string    `json:"filename"`
}

type Comment struct {
	Comment_id int       `json:"comment_id"`
	Post_id    int       `json:"post_id"`
	CreatedAt  time.Time `json:"created_at"`
	Likes      int       `json:"likes"`
	Dislikes   int       `json:"dislikes"`
	Content    string    `json:"content"`
	UserLiked  int       `json:"UserLiked"`
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
