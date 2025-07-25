package v1

import (
	v1 "api-gateway/docs/proto/v1"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	pbgrpc "google.golang.org/grpc"
)

func NewRoutes(g *echo.Group, authClient *pbgrpc.ClientConn, l *zap.SugaredLogger) {
	api := V1{
		authClient: v1.NewAuthClient(authClient),
		l: l,
	}

	authGroup := g.Group("/auth")
	{
		authGroup.POST("/login", api.Login)
		authGroup.POST("/register", api.Register)
	}

	//userGroup := g.Group("/users", )
}
