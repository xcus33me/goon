package service

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"goon/services/auth/internal/entity"
	"goon/services/auth/internal/repo"
	e "goon/services/auth/pkg/errors"
)

var (
	MaxPasswordLength = 255
)

type AuthService struct {
	repo repo.AccountRepo
}

func NewAuthService(repo repo.AccountRepo) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Login(login, password string) (*entity.Account, error) {
	account, err := s.repo.FindByLogin(login)
	if err != nil {
		return nil, e.UserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, e.WrongCredentials
		}
		return nil, fmt.Errorf("%w: %v", e.HashingFailed, err)
	}

	return account, nil
}

func (s *AuthService) Register(login, password string) (*entity.Account, error) {
	_, err := s.repo.FindByLogin(login)
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

	account := &entity.Account{
		Login:        login,
		PasswordHash: string(hash),
	}

	if err := s.repo.CreateAccount(account); err != nil {
		return nil, err
	}

	return account, nil
}
