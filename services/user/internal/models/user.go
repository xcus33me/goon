package models

import (
	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Bio       string    `json:"bio"`
}
