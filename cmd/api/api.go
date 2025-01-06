package api

import (
	"goon/service/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type APIServer struct {
	addr string
	db   *gorm.DB
}

func NewAPIServer(addr string, db *gorm.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := gin.Default()
	userStore := user.NewStore(s.db)

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			userHandler := user.NewHandler(userStore)
			userHandler.RegisterRoutes(v1.Group("/users"))
		}
	}

	return router.Run(s.addr)
}
