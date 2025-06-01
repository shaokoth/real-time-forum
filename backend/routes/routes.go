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

	mux.HandleFunc("/ws", handlers.AuthMiddleware(http.HandlerFunc(handlers.HandleWebSocket)))
	mux.HandleFunc("/posts", handlers.AuthMiddleware(http.HandlerFunc(handlers.HandlePosts)))
	mux.HandleFunc("/users", handlers.AuthMiddleware(http.HandlerFunc(handlers.HandleUsers)))
	mux.HandleFunc("/comments", handlers.AuthMiddleware(http.HandlerFunc(handlers.HandleComments)))
	mux.HandleFunc("/categories", handlers.AuthMiddleware(http.HandlerFunc(handlers.HandleGetCategories)))

	return mux, nil
}
