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

	return mux, nil
}
