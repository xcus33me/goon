package app

import (
	"user/config"
	"user/pkg/logger"
	"user/pkg/postgres"
)

func Run(cfg *config.Config) error {
	l := logger.New(cfg.Log.Level)

	pg, err := postgres.New(cfg)
}
