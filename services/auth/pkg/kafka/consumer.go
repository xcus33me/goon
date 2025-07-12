package kafka

import (
	"auth/pkg/logger"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	logger logger.Interface
}
