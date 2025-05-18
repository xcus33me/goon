package persistent

import (
	"auth/internal/entity"
	"auth/pkg/postgres"
)

type AuthRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *AuthRepo {
	return &AuthRepo{pg}
}

func (r *AuthRepo) CreateUser(user *entity.User) error {
	query := "INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id"
	return r.Db.QueryRow(query, user.Login, user.PasswordHash).Scan(&user.ID)
}

func (r *AuthRepo) FindByLogin(login string) (*entity.User, error) {
	query := "SELECT id, login, password_hash FROM users WHERE login = $1"
	user := &entity.User{}

	err := r.Db.QueryRow(query, login).Scan(&user.ID, &user.Login, &user.PasswordHash)
	if err != nil {
		return nil, err
	}

	return user, err
}
