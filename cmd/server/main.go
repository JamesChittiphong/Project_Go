package main

import (
	"log"

	"Backend_Go/internal/app"
	"Backend_Go/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env.config")

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	app := app.NewApp(db)

	log.Println("Server running on :8081")
	app.Listen(":8081")
}
