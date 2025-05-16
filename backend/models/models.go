package models

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

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

type Categories struct {
	All_Categories string
	Technology     string
	Health         string
	Math           string
	Nature         string
	Science        string
	Religion       string
	Education      string
	Politics       string
	Fashion        string
	Lifestyle      string
	Sports         string
	Arts           string
}

type Users struct {
	Username string
	Email    string
	Password string
}


type Comment struct {
	Comment_id int       `json:"comment_id"`
	Post_id    string    `json:"post_id"`
	CreatedAt  time.Time `json:"created_at"`
	Likes      int       `json:"likes"`
	Dislikes   int       `json:"dislikes"`
	Content    string    `json:"content"`
}

type Image struct {
	ImageID   string    `json:"image_id"`
	UserID    string    `json:"user_id"`
	PostID    string    `json:"post_id"`
	Filename  string    `json:"filename"`
	Path      string    `json:"path"`
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
