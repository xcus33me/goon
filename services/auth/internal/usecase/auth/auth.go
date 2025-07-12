package auth

import (
	"auth/internal/entity"
	"auth/internal/repo"
	e "auth/pkg/errors"
	"auth/pkg/kafka"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	repo repo.AuthRepo
	kafkaProducer *kafka.Producer
	jwtSecret string
}

func New(repo repo.AuthRepo, kafkaProducer *kafka.Producer, jwtSecret string) *UseCase {
	return &UseCase{
		repo:      repo,
		kafkaProducer: kafkaProducer,
		jwtSecret: jwtSecret,
	}
}

func (uc *UseCase) generateToken(user *entity.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":   user.ID,
		"login": user.Login,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("%w: could not sign token: %v", e.FailedToGenerateToken, err)
	}

	return tokenString, nil
}

func (uc *UseCase) Login(ctx context.Context, login, password string) (*entity.User, string, error) {
	user, err := uc.repo.FindByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return nil, "", e.WrongCredentials
		}
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, "", e.WrongCredentials
		}
		return nil, "", fmt.Errorf("%w: %v", e.HashingFailed, err)
	}

	token, err := uc.generateToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %v", e.FailedToGenerateToken, err)
	}

	return user, token, nil
}

func (uc *UseCase) Register(ctx context.Context, login, password string) (*entity.User, error) {
	_, err := uc.repo.FindByLogin(ctx, login)
	if err == nil {
		return nil, e.UserAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", e.HashingFailed, err)
	}

	user := &entity.User{
		Login:        login,
		PasswordHash: string(hash),
	}

	if err := uc.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// kafka
	go func() {
		eventCtx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()


	}

	return user, nil
}

func (uc *UseCase) UpdatePassword(ctx context.Context, userID int64, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%w: %v", e.HashingFailed, err)
	}

	request := &repo.UpdatePasswordRequest{
		ID:           userID,
		PasswordHash: string(hash),
	}

	_, err = uc.repo.UpdatePasswordByID(ctx, request)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return e.UserNotFound
		}
		return err
	}

	return nil
}
