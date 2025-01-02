package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

type User struct {
	gorm.Model
	ID                uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name              string    `gorm:"type:varchar(100);not null"`
	Email             string    `gorm:"type:varchar(255);unique;not null"`
	Verified          bool      `gorm:"default:false;not null"`
	Password          string    `gorm:"type:varchar(100);not null"`
	VerificationToken string    `gorm:"type:varchar(255)"`
	Role              UserRole  `gorm:"type:varchar(5); default:'user';not null"`
	TokenExpiresAt    time.Time
}
