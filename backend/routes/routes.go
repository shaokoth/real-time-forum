package routes

import (
	"net/http"

	"real-time-forum/backend/handlers"
)

func Routers() (*http.ServeMux, error) {
	mux := http.NewServeMux()

	//  Serve static files
	fileServer := http.FileServer(http.Dir("frontend/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	scriptServer := http.FileServer(http.Dir("frontend/js/"))
	mux.Handle("/js/", http.StripPrefix("/js/", scriptServer))

	mux.HandleFunc("/register", handlers.RegisterUser)
	mux.HandleFunc("/login", handlers.HandleLogin)
	mux.HandleFunc("/logout", handlers.LogoutUser)

	mux.HandleFunc("/ws", handlers.AuthMiddleware(http.HandlerFunc(handlers.HandleWebSocket)))
	mux.HandleFunc("/posts", handlers.AuthMiddleware(http.HandlerFunc(handlers.HandlePosts)))
	mux.HandleFunc("/users", handlers.AuthMiddleware(http.HandlerFunc(handlers.HandleUsers)))
	mux.HandleFunc("/comments", handlers.AuthMiddleware(http.HandlerFunc(handlers.HandleComments)))

	return mux, nil
}
