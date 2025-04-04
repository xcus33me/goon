package models

type Account struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
