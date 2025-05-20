package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"real-time-forum/backend/database"
	"real-time-forum/backend/routes"
)

func main() {
	database.Init()
	defer database.Db.Close()

	mux, err := routes.Routers()
	if err != nil {
		fmt.Println("Error")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "1995"
	}

	log.Println("server started on port http://localhost:1995")
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Println("Error starting server")
	}
}
