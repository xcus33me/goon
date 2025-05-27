package auth

import (
	"auth/internal/entity"
	"auth/internal/repo"
	e "auth/pkg/errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UseCase struct {
	repo repo.AuthRepo
	//webAPI    repo.AuthWebAPI
	jwtSecret string
}

func New(repo repo.AuthRepo, jwtSecret string) *UseCase {
	return &UseCase{
		repo:      repo,
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
		return "", err
	}

	return tokenString, nil
}

func (uc *UseCase) Login(login, password string) (*entity.User, string, error) {
	user, err := uc.repo.FindByLogin(login)
	if err != nil {
		return nil, "", e.UserNotFound
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

func (uc *UseCase) Register(login, password string) (*entity.User, error) {
	_, err := uc.repo.FindByLogin(login)
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

	if err := uc.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}
