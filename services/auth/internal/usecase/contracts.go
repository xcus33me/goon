package usecase

import (
	"auth/internal/entity"
	"context"
)

//go:generate mockgen -source=./internal/usecase/contracts.go -destination=./internal/usecase/mocks_usecase_test.go -package=usecase_test

type (
	Auth interface {
		Login(ctx context.Context, login, password string) (*entity.User, string, error)
		Register(ctx context.Context, login, password string) (*entity.User, error)
		UpdatePassword(ctx context.Context, userID int64, password string) error
	}
)
