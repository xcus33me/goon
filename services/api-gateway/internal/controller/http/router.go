package http

import (
	v1 "api-gateway/internal/controller/http/v1"

	pbgrpc "google.golang.org/grpc"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func NewRouter(e *echo.Echo, authClient *pbgrpc.ClientConn, l *zap.SugaredLogger) {
	appV1Group := e.Group("/v1")
	{
		v1.NewRoutes(appV1Group, authClient, l)
	}
}
