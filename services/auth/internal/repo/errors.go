package repo

import "errors"

var (
	ErrNotFound       = errors.New("entity not found")
	ErrDuplicateEntry = errors.New("entity already exists")
)
