package repo

import (
	"auth/internal/entity"
	"auth/pkg/postgres"
)

type UserRepo interface {
	CreateUser(user *entity.User) error
	FindByLogin(login string) (*entity.User, error)
}

type userRepo struct {
	*postgres.Postgres
}

func NewUserStorage(pg *postgres.Postgres) *userRepo {
	return &userRepo{pg}
}

func (s *userRepo) CreateUser(user *entity.User) error {
	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id"
	return s.Db.QueryRow(query, user.Login, user.PasswordHash).Scan(&user.ID)
}

func (s *userRepo) FindByLogin(login string) (*entity.User, error) {
	query := "SELECT id, login, password_hash FROM users WHERE login = $1"
	user := &entity.User{}

	err := s.Db.QueryRow(query, login).Scan(&user.ID, &user.Login, &user.PasswordHash)
	if err != nil {
		return nil, err
	}

	return user, err
}
