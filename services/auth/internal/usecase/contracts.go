package usecase

import "auth/internal/entity"

//go:generate mockgen -source=./internal/usecase/contracts.go -destination=./internal/usecase/mocks_usecase_test.go -package=usecase_test

type (
	Auth interface {
		Login(login, password string) (*entity.User, string, error)
		Register(login, password string) (*entity.User, error)
	}
)
