package repo

import "auth/internal/entity"

type (
	AuthRepo interface {
		CreateUser(user *entity.User) error
		FindByLogin(login string) (*entity.User, error)
	}

	AuthWebAPI interface {
		Notify()
	}
)
