package routes

import (
	"net/http"

	"real-time-forum/backend/handlers"
)

func Routers() (*http.ServeMux, error) {
	mux := http.NewServeMux()

	//  Serve static files
	fileServer := http.FileServer(http.Dir("frontend/css"))
	mux.Handle("/css/", http.StripPrefix("/css/", fileServer))
	scriptServer := http.FileServer(http.Dir("frontend/js/"))
	mux.Handle("/js/", http.StripPrefix("/js/", scriptServer))
	imageServer := http.FileServer(http.Dir("frontend/image"))
	mux.Handle("/static/image/", http.StripPrefix("/static/image/", imageServer))
    
	mux.HandleFunc("/", handlers.HandleHomepage)
	mux.HandleFunc("/register", handlers.RegisterUser)
	mux.HandleFunc("/login", handlers.HandleLogin)
	mux.HandleFunc("/logout", handlers.LogoutUser)

	// Posts and categories can be viewed without authentication
	mux.HandleFunc("/posts", handlers.HandlePosts)
	mux.HandleFunc("/comments", handlers.HandleComments)
	mux.HandleFunc("/categories", handlers.HandleGetCategories)
	
	// These endpoints require authentication
	mux.HandleFunc("/ws", handlers.AuthMiddleware(http.HandlerFunc(handlers.HandleWebSocket)))
	mux.HandleFunc("/users", handlers.AuthMiddleware(http.HandlerFunc(handlers.HandleUsers)))
	mux.HandleFunc("/posts/like", handlers.AuthMiddleware(http.HandlerFunc(handlers.LikePostHandler)))
	mux.HandleFunc("/posts/dislike", handlers.AuthMiddleware(http.HandlerFunc(handlers.DislikePostHandler)))
	mux.HandleFunc("/comments/like", handlers.AuthMiddleware(http.HandlerFunc(handlers.LikeCommentHandler)))
	mux.HandleFunc("/comments/dislike", handlers.AuthMiddleware(http.HandlerFunc(handlers.DislikeCommentHandler)))

	return mux, nil
}
