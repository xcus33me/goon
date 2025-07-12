package kafka

import (
	"auth/config"
	"auth/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type ProducerInterface interface {
	PublishUserCreated(ctx context.Context, userID int64) error
	PublishUserUpdated(ctx context.Context, userID int64, time time.Time) error
}

type Producer struct {
	writer *kafka.Writer
	logger logger.Interface
}

func NewProducer(cfg *config.Kafka, l logger.Interface) *Producer {
	writer := kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Balancer: &kafka.Hash{},

		RequiredAcks: kafka.RequiredAcks(cfg.Producer.RequireAcks),

		WriteTimeout: cfg.Producer.WriteTimeout,
		ReadTimeout:  cfg.Producer.ReadTimeout,

		AllowAutoTopicCreation: cfg.Producer.AllowAutoTopicCreation,
	}

	return &Producer{
		writer: &writer,
		logger: l,
	}
}

func (p *Producer) PublishUserCreated(ctx context.Context, userID int64) error {
	event := UserCreated{
		UserID:    userID,
		EventID:   uuid.New().String(),
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("kafka - PublishUserCreate - failed to marshal event: %w", err)
	}

	message := kafka.Message{
		Topic: "user.created",
		Key:   []byte(fmt.Sprintf("user_%d", userID)),
		Value: data,
		Time:  time.Now(),
	}

	if err := p.writer.WriteMessages(ctx, message); err != nil {
		return fmt.Errorf("kafka - PublishUserCreate - failed to write messages: %w", err)
	}

	p.logger.Info("kafka - PublishUserCreated: published")

	return nil
}
