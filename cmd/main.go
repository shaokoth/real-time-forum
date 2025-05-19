package main

import (
	"real-time-forum/backend/database"
)

func main() {
	database.Init()
	defer database.Db.Close()
}
