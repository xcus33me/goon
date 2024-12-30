package main

import (
	database "goon/db"
	repo "goon/repository"
	"log"
	"time"

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

	user, err := repo.CreateUser("block", "ngachain@yandex.ru", "12345", "123", time.Now())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	println(user)
}
