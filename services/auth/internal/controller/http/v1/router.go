package v1

import (
	"auth/internal/usecase"
	"auth/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func NewAuthRoutes(g *echo.Group, a usecase.Auth, l logger.Interface) {
	api := &V1{
		a: a,
		l: l,
		v: validator.New(validator.WithRequiredStructEnabled()),
	}

	authGroup := g.Group("/auth")
	{
		authGroup.POST("/login", api.Login)
		authGroup.POST("/register", api.Register)
	}
}
