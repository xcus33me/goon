package database

import (
	"goon/data"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	dsn := os.Getenv("DB_URL")

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	return nil
}

func Sync() {
	DB.AutoMigrate(&data.User{})
}

func CleanUp() {
	DB.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
}
