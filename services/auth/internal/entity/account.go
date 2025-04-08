package entity

type Account struct {
	ID           int    `json:"id"`
	Login        string `json:"name"`
	PasswordHash string `json:"password_hash"`
}
