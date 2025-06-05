package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"real-time-forum/backend/database"
	"real-time-forum/backend/handlers"
	"real-time-forum/backend/routes"
)

func main() {
	database.Init()
	handlers.InitWebSocket()

	defer database.Db.Close()

	mux, err := routes.Routers()
	if err != nil {
		fmt.Println("Error")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("server started on port http://localhost:8080")
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Println("Error starting server")
	}
}
