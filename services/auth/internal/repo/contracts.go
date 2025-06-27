package repo

import (
	"auth/internal/entity"
	"context"
	"time"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	AuthRepo interface {
		CreateUser(ctx context.Context, user *entity.User) error
		FindByLogin(ctx context.Context, login string) (*entity.User, error)
		UpdatePasswordByID(ctx context.Context, request *UpdatePasswordRequest) (*UpdatePasswordResponse, error)
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
