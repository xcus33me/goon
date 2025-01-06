package db

import (
	"fmt"
	"goon/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func createPostgresDSN(cfg config.Config) string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBAddress,
		cfg.DBName,
	)
}

func NewPostgresStorage(cfg config.Config) (*gorm.DB, error) {
	dsn := createPostgresDSN(cfg)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)

	}

	return db, nil
}
