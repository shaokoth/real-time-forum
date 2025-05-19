package main

import (
	"fmt"
	"log"
	"net/http"

	"real-time-forum/backend/database"
	"real-time-forum/backend/routes"
)

func main() {
	database.Init()
	defer database.Db.Close()

	// Get router ServerMux
	mux, err := routes.Routers()
	if err != nil {
		fmt.Errorf("error intialize routes: %v", err)
		return
	}
	// start HTTP server
	log.Println("starting server on http://localhost:8000")
	err = http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Println("Server failed: %v", err)
		return
	}
}
