package v1

import (
	"net/http"

	v1 "github.com/xcus33me/protos/gen/go/auth"

	"github.com/labstack/echo/v4"
)

type (
	loginRequest struct {
		Login string `json:"login"`
		Password string `json:"password"`
	}

	registerRequest struct {
		Login string `json:"login"`
		Password string `json:"password"`
	}
)

func (api *V1) Login(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return api.errorResponse(c, http.StatusBadRequest, err)
	}

	res, err := api.authClient.Login(c.Request().Context(), &v1.LoginRequest{
		Login: req.Login,
		Password: req.Password,
	})
	if err != nil {
		api.l.Errorf("v1 - Login: grpc call to auth service failed: %v", err)
		api.errorResponse(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, res)
}

func (api *V1) Register(c echo.Context) error {
	var req registerRequest
	if err := c.Bind(&req); err != nil {
		return api.errorResponse(c, http.StatusBadRequest, err)
	}

	res, err := api.authClient.Register(c.Request().Context(), &v1.RegisterRequest{
		Login: req.Login,
		Password: req.Password,
	})
	if err != nil {
		api.l.Errorf("v1 - Register: grpc call to auth service failed: %v", err)
		api.errorResponse(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, res)
}
