package app

import (
	"auth/config"
	"auth/internal/controller/http"
	"auth/internal/repo/persistent"
	"auth/internal/usecase/auth"
	"auth/pkg/httpserver"
	"auth/pkg/kafka"
	"auth/pkg/logger"
	"auth/pkg/postgres"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	// Logger
	l := logger.NewWithFile(cfg.Log.Level, cfg.Log.Path)

	// Postgres
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %v", err))
	}
	defer pg.Close()

	// Kafka Producer
	producer := kafka.NewProducer(&cfg.Kafka, l)

	l.Info("Auth service in action🔥🚀")

	// Usecase
	authUseCase := auth.New(
		persistent.New(pg),
		producer,
		cfg.Auth.JWTSecret,
	)

	// HttpServer
	httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	http.NewRouter(httpServer.App, cfg, authUseCase, l)

	// Start servers
	httpServer.Run()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
