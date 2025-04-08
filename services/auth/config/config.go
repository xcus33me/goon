package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	App App
}

type App struct {
	Name    string `env:"APP_NAME,required"`
	Version string `env:"APP_VERSION,required"`
}

type Log struct {
	Level string `env:"LOG_LEVEL,required"`
}

type PG struct {
	PoolMax int    `env:"PG_POOL_MAX,required"`
	URL     string `env:"PG_URL,required"`
}

type RMQ struct {
	ServerExchange string `env:"RMQ_RPC_SERVER,required"`
	ClientExchange string `env:"RMQ_RPC_CLIENT,required"`
	URL            string `env:"RMQ_URL,required"`
}

type Metrics struct {
	Enabled bool `env:"METRICS_ENABLED" envDefault:"true"`
}

type Swagger struct {
	Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
