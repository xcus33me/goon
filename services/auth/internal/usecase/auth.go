package usecase

import (
	"auth/internal/entity"
	"auth/internal/repo"
	e "auth/pkg/errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	MaxLoginLength    = 30
	MinLoginLength    = 4
	MaxPasswordLength = 255
	MinPasswordLength = 6
)

type AuthUseCase struct {
	repo repo.UserRepo
}

func NewAuthUseCase(repo repo.UserRepo) *AuthUseCase {
	return &AuthUseCase{repo: repo}
}

func (uc *AuthUseCase) Login(login, password string) (*entity.User, error) {
	user, err := uc.repo.FindByLogin(login)
	if err != nil {
		return nil, e.UserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, e.WrongCredentials
		}
		return nil, fmt.Errorf("%w: %v", e.HashingFailed, err)
	}

	return user, nil
}

func (uc *AuthUseCase) Register(login, password string) (*entity.User, error) {
	_, err := uc.repo.FindByLogin(login)
	if err == nil {
		return nil, e.UserAlreadyExists
	}

	if len(password) > MaxPasswordLength {
		return nil, fmt.Errorf("%w: %v", e.ExceededMaxPasswordLength, MaxPasswordLength)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", e.HashingFailed, err)
	}

	user := &entity.User{
		Login:        login,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}
