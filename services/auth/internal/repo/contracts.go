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
		UpdatePasswordByID(request *UpdatePasswordRequest) (*UpdatePasswordResponse, error)
	}

	AuthWebAPI interface {
		Notify()
	}
)

type UpdatePasswordRequest struct {
	ID           int64
	PasswordHash string
}

type UpdatePasswordResponse struct {
	ID        int64
	UpdatedAt time.Time
}
