package http

import (
	"auth/config"
	m "auth/internal/controller/http/middleware"
	v1 "auth/internal/controller/http/v1"
	"auth/internal/usecase"
	"auth/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(e *echo.Echo, cfg *config.Config, a usecase.Auth, l logger.Interface) {
	e.Use(m.Logger(l))
	e.Use(middleware.Recover())

	//if cfg.Metrics.Enabled {
	//	e.GET("/swagger/*", echoSwagger.WrapHandler)
	//}

	appV1Group := e.Group("/v1")
	{
		v1.NewAuthRoutes(appV1Group, a, l)
	}
}
