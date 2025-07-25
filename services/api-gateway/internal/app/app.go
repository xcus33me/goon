package app

import (
	"api-gateway/config"
	"api-gateway/internal/controller/http"
	"api-gateway/pkg/grpcclient"
	"api-gateway/pkg/httpserver"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func Run(cfg *config.Config) {
	zaplogger, err := zap.NewProduction()
	if err != nil {
		panic("Failed to create logger")
	}
	l := zaplogger.Sugar()

	// gRPC Clients
	authClient, err := grpcclient.New(cfg.GRPC.TargetAuth)
	if err != nil {
		l.Fatalf("Failed to connect to auth GRPC Server: %w", err)
	}
	l.Infof("✅ Successfully connected to auth service at %s", cfg.GRPC.TargetAuth)


	// HTTP Server
	httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	http.NewRouter(httpServer.App, authClient.Conn, l)

	// Start Servers
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
	authClient.Close()
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
