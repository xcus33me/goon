package http

import (
	"auth/config"
	"auth/internal/usecase"
	"auth/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func NewRouter(e *echo.Echo, cfg *config.Config, l logger.Logger, t usecase.AuthUseCase) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	if cfg.Metrics.Enabled {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	appV1Group := e.Group("v1")
	{

	}
}
