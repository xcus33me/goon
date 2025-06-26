package app

import (
	"auth/pkg/logger"
	"errors"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	_defaultAttempts = 3
	_defaultTimeout  = time.Second
)

func init() {
	l := logger.NewWithFile("info", os.Getenv("LOG_PATH"))

	databaseURL, ok := os.LookupEnv("PG_URL")
	if !ok || len(databaseURL) == 0 {
		l.Fatal("migrate: environment variable not declared: PG_URL")
	}

	databaseURL += "?sslmode=disable"

	var (
		attempts = _defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://migrations", databaseURL)
		if err == nil {
			break
		}

		l.Info("migrate: postgres is trying to connect, attempts left: ", attempts)
		time.Sleep(_defaultTimeout)
		attempts--
	}

	if err != nil {
		l.Fatal("migrate: postgres connect error: ", err)
	}

	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		l.Fatal("migrate: up error: ", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		l.Info("migrate: no change")
		return
	}

	l.Info("Migrate: up success")
}
