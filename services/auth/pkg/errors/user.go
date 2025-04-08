package errors

import "errors"

var (
	UserAlreadyExists         = errors.New("user already exists")
	UserNotFound              = errors.New("account not found")
	HashingFailed             = errors.New("hashing failed")
	WrongCredentials          = errors.New("wrong credentials")
	ExceededMaxPasswordLength = errors.New("exceeded max password length")
)
