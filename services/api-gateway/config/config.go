package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	App  App
	HTTP HTTP
	GRPC GRPC
	Log  Log
	//Metrics Metrics
	//Swagger Swagger
}

type App struct {
	Name    string `env:"APP_NAME,required"`
	Version string `env:"APP_VERSION,required"`
}

type HTTP struct {
	Port           string `env:"HTTP_PORT,required"`
	UsePreforkMode bool   `env:"HTTP_USE_PREFORK_MODE" envDefault:"false"`
}

type GRPC struct {
	TargetAuth string `env:"GRPC_TARGET_AUTH,required"`
	TargetUser string `env:"GRPC_TARGET_USER,required"`
}

type Log struct {
	Level string `env:"LOG_LEVEL,required"`
	Path  string `env:"LOG_PATH,required"`
}

// type Metrics struct {
// 	Enabled bool `env:"METRICS_ENABLED" envDefault:"true"`
// }

// type Swagger struct {
// 	Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
// }

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
