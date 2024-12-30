package repo

import (
	"goon/data"
	database "goon/db"
	"goon/utils"
	"time"

	"github.com/google/uuid"
)

func CreateUser(name, email, password, verificationToken string, tokenExpiresAt time.Time) (*data.User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &data.User{
		ID:                uuid.New(),
		Name:              name,
		Email:             email,
		Verified:          false,
		Password:          hashedPassword,
		VerificationToken: verificationToken,
		Role:              data.RoleUser,
		TokenExpiresAt:    tokenExpiresAt,
	}

	if err := database.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
