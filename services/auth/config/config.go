package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	App  App
	HTTP HTTP
	Log  Log
	PG   PG
	//RMQ     RMQ
	Metrics Metrics
	Swagger Swagger
	Auth    Auth
}

type App struct {
	Name    string `env:"APP_NAME,required"`
	Version string `env:"APP_VERSION,required"`
}

type HTTP struct {
	Port           string `env:"HTTP_PORT,required"`
	UsePreforkMode bool   `env:"HTTP_USE_PREFORK_MODE" envDefault:"false"`
}

type Log struct {
	Level string `env:"LOG_LEVEL,required"`
	Path  string `env:"LOG_PATH,required"`
}

type PG struct {
	PoolMax int    `env:"PG_POOL_MAX,required"`
	URL     string `env:"PG_URL,required"`
}

type Kafka struct {
	Brokers  []string `env:"KAFKA_BROKERS,required" envSeparator:","`
	Producer ProducerConfig
	Conumer  ConsumerConfig
}

type ProducerConfig struct {
	// NoResponse(0), WaitForLocal(1), WaitForAll(-1)
	RequireAcks int `env:"KAFKA_PRODUCER_ACKS" envDefault:"-1"`

	// Retry
	MaxRetries   int           `env:"KAFKA_PRODUCER_MAX_RETRIES" envDefault:"3"`
	RetryBackoff time.Duration `env:"KAFKA_PRODUCER_RETRY_BACKOFF" envDefault:"100ms"`

	// Timeouts
	WriteTimeout time.Duration `env:"KAFKA_PRODUCER_WRITE_TIMEOUT" envDefault:"10s"`
	ReadTimeout  time.Duration `env:"KAFKA_PRODUCER_READ_TIMEOUT" envDefault:"10s"`

	// Batching
	// BatchSize    int           `env:"KAFKA_PRODUCER_BATCH_SIZE" envDefault:"16384"` // 16KB
	// BatchTimeout time.Duration `env:"KAFKA_PRODUCER_BATCH_TIMEOUT" envDefault:"10ms"`
	// BatchBytes   int           `env:"KAFKA_PRODUCER_BATCH_BYTES" envDefault:"1048576"` // 1MB

	// Connection pooling
	// MaxOpenRequests int `env:"KAFKA_PRODUCER_MAX_OPEN_REQUESTS" envDefault:"5"`

	// Auto topic creation
	AllowAutoTopicCreation bool `env:"KAFKA_PRODUCER_AUTO_TOPIC_CREATION" envDefault:"false"`
}

type ConsumerConfig struct {
	GroupID string `env:"KAFKA_CONSUMER_GROUP_ID,required"`

	// Batch processing
	// MinBytes        int           `env:"KAFKA_CONSUMER_MIN_BYTES" envDefault:"1"`
	// MaxBytes        int           `env:"KAFKA_CONSUMER_MAX_BYTES" envDefault:"1048576"`       // 1MB
	// MaxWait         time.Duration `env:"KAFKA_CONSUMER_MAX_WAIT" envDefault:"500ms"`

	// Commit
	CommitInterval    time.Duration `env:"KAFKA_CONSUMER_COMMIT_INTERVAL" envDefault:"1s"`
	Partition         int           `env:"KAFKA_CONSUMER_PARTITION" envDefault:"-1"` // -1 = all partitions
	HeartbeatInterval time.Duration `env:"KAFKA_CONSUMER_HEARTBEAT_INTERVAL" envDefault:"3s"`
	SessionTimeout    time.Duration `env:"KAFKA_CONSUMER_SESSION_TIMEOUT" envDefault:"30s"`

	// Retry
	MaxRetries   int           `env:"KAFKA_CONSUMER_MAX_RETRIES" envDefault:"3"`
	RetryBackoff time.Duration `env:"KAFKA_CONSUMER_RETRY_BACKOFF" envDefault:"250ms"`
}

type Metrics struct {
	Enabled bool `env:"METRICS_ENABLED" envDefault:"true"`
}

type Swagger struct {
	Enabled bool `env:"SWAGGER_ENABLED" envDefault:"false"`
}

type Auth struct {
	JWTSecret string `env:"JWT_SECRET,required"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
