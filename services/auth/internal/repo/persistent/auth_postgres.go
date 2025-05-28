package persistent

import (
	"auth/internal/entity"
	"auth/internal/repo"
	"auth/pkg/postgres"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

type AuthRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *AuthRepo {
	return &AuthRepo{pg}
}

func (r *AuthRepo) CreateUser(user *entity.User) error {
	query := "INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id, created_at, updated_at"
	err := r.Db.QueryRow(query, user.Login, user.PasswordHash).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return repo.ErrDuplicateEntry
			}
		}
		return fmt.Errorf("auth_postgres - CreateUser: failed to create user: %w", err)
	}
	return nil
}

func (r *AuthRepo) FindByLogin(login string) (*entity.User, error) {
	query := "SELECT id, login, password_hash FROM users WHERE login = $1"
	user := &entity.User{}

	err := r.Db.QueryRow(query, login).Scan(&user.ID, &user.Login, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repo.ErrNotFound
		}
		return nil, fmt.Errorf("auth_postgres - FindByLogin: failed to find user: %w", err)
	}

	return user, err
}

func (r *AuthRepo) UpdatePasswordByID(request *repo.UpdatePasswordRequest) (*repo.UpdatePasswordResponse, error) {
	query := `
		UPDATE users SET password_hash = $1, updated_at = Now() 
		WHERE id = $2
		RETURNING id, updated_at`

	response := &repo.UpdatePasswordResponse{}

	err := r.Db.QueryRow(query, request.PasswordHash, request.ID).Scan(&response.ID, &response.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repo.ErrNotFound
		}
		return nil, fmt.Errorf("auth_postgres - UpdatePasswordByID: failed to update password: %w", err)
	}

	return response, nil
}
