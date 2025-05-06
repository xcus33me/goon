package usecase

import "auth/internal/entity"

type (
	Auth interface {
		Login(login, password string) (*entity.User, string, error)
		Register(login, password string) (*entity.User, error)
	}
)
