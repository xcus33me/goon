package main

import (
	"fmt"
	"goon/cmd/api"
	"goon/config"
	"goon/db"
	"log"

	"gorm.io/gorm"
)

func main() {
	cfg := config.Envs

	db, err := db.NewPostgresStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	InitStorage(db)

	server := api.NewAPIServer(fmt.Sprintf(":%s", config.Envs.Port), db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func InitStorage(db *gorm.DB) {
	sqlDb, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	err = sqlDb.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: ✅ Successfully connected!")
}
