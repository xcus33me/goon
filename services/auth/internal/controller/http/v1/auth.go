package v1

import (
	"auth/internal/usecase"
	"auth/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type authRoutes struct {
	t usecase.AuthUseCase
	l logger.Logger
	v *validator.Validate
}

func NewAuthRoutes(g *echo.Group, t usecase.AuthUseCase, l logger.Logger) {
	r := &authRoutes{
		t: t,
		l: l,
		v: validator.New(),
	}
}
