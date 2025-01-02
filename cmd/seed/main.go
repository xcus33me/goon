package main

import (
	database "goon/db"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	err := database.Connect()
	if err != nil {
		log.Printf("❌ Failed to connect to database: %v", err)
	}
	log.Println("✅ Successfully connected to the database!")

	database.CleanUp()
	database.Sync()

	// todo: Register() here.
}
