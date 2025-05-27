package repo

import (
	"auth/internal/entity"
	"time"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	AuthRepo interface {
		CreateUser(user *entity.User) error
		FindByLogin(login string) (*entity.User, error)
	}

	AuthWebAPI interface {
		Notify()
	}
)

type UpdatePasswordDTO struct {
	ID              int64
	NewPasswordHash string
	UpdatedAt       time.Time
}
