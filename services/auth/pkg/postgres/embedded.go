package postgres

import (
	"fmt"
	"sync/atomic"
	"testing"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Embedded struct {
	postgres *embeddedpostgres.EmbeddedPostgres
	DB       *sqlx.DB
	Port     uint32
}

var portCounter uint32 = 15000

// NewEmbeddedDB создает embedded PostgreSQL с автоинкрементом порта
func NewEmbedded(t *testing.T) *Embedded {
	t.Helper()

	port := atomic.AddUint32(&portCounter, 1)

	postgres := embeddedpostgres.NewDatabase(
		embeddedpostgres.DefaultConfig().
			Port(port).
			Database("testdb").
			Username("testuser").
			Password("testpass"),
	)

	err := postgres.Start()
	if err != nil {
		t.Fatalf("Failed to start embedded postgres: %v", err)
	}

	dsn := fmt.Sprintf("host=localhost port=%d user=testuser password=testpass dbname=testdb sslmode=disable", port)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		postgres.Stop()
		t.Fatalf("Failed to connect: %v", err)
	}

	// Создаем схему
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			login VARCHAR(30) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		postgres.Stop()
		t.Fatalf("Failed to create schema: %v", err)
	}

	embeddedDB := &Embedded{
		postgres: postgres,
		DB:       db,
		Port:     port,
	}

	t.Cleanup(func() {
		embeddedDB.Close()
	})

	return embeddedDB
}

func (e *Embedded) Close() {
	if e.DB != nil {
		e.DB.Close()
	}
	if e.postgres != nil {
		e.postgres.Stop()
	}
}
