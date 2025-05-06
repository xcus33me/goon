package v1

import (
	"auth/internal/usecase"
	"auth/pkg/logger"

	"github.com/go-playground/validator/v10"
)

type V1 struct {
	a usecase.Auth
	l logger.Interface
	v *validator.Validate
}
