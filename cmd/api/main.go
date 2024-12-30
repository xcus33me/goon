package main

import (
	database "goon/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

	database.Sync()

	router := gin.Default()
	router.GET("/", homepage)
	router.GET("/cat", getCat)

	err = router.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func homepage(c *gin.Context) {
	c.String(http.StatusOK, "Welcome!")
}

func getCat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "mew ^3^",
	})
}
