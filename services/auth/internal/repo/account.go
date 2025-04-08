package repo

import (
	"database/sql"
	"goon/services/auth/internal/entity"
)

type AccountRepo interface {
	CreateAccount(account *entity.Account) error
	FindByLogin(login string) (*entity.Account, error)
}

type AccountStorage struct {
	db *sql.DB
}

func NewAccountStorage(db *sql.DB) *AccountStorage {
	return &AccountStorage{db: db}
}

func (s *AccountStorage) CreateAccount(account *entity.Account) error {
	query := "INSERT INTO accounts (username, password_hash) VALUES ($1, $2) RETURNING id"
	return s.db.QueryRow(query, account.Login, account.PasswordHash).Scan(&account.ID)
}

func (s *AccountStorage) FindByLogin(login string) (*entity.Account, error) {
	query := "SELECT id, login, password_hash FROM accounts WHERE login = $1"
	account := &entity.Account{}

	err := s.db.QueryRow(query, login).Scan(&account.ID, &account.Login, &account.PasswordHash)
	if err != nil {
		return nil, err
	}

	return account, err
}
