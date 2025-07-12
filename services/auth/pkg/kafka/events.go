package kafka

import "time"

type (
	UserCreated struct {
		UserID    int64     `json:"user_id"`
		EventID   string    `json:"event_id"`
		Timestamp time.Time `json:"timestamp"`
	}
)
