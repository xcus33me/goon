package repo

import (
	database "goon/db"
	"goon/models"
	"time"

	"github.com/google/uuid"
)

func SaveUser(name, email, password, verificationToken string, expiresAt time.Time) (*models.User, error) {
	user := models.User{
		ID:                uuid.New(),
		Name:              name,
		Email:             email,
		Verified:          false,
		Password:          password,
		VerificationToken: verificationToken,
		Role:              models.RoleUser,
		TokenExpiresAt:    expiresAt,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
