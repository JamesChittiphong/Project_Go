package main

import (
	"log"
	"os"

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

	// Serve static files from /uploads
	// Ensure the directory exists
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		log.Println("Warning: ./uploads directory does not exist")
	}
	app.Static("/uploads", "./uploads")

	log.Println("Server running on :8081")
	app.Listen(":8081")
}
