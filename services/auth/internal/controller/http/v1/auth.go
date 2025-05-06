package v1

import (
	"auth/internal/controller/http/v1/request"
	"auth/internal/controller/http/v1/response"
	e "auth/pkg/errors"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (api *V1) Login(c echo.Context) error {
	var body request.Login
	if err := c.Bind(&body); err != nil {
		api.errorResponse(c, http.StatusBadRequest, err)
	}

	if err := api.v.Struct(&body); err != nil {
		api.errorResponse(c, http.StatusBadRequest, err)
	}

	user, token, err := api.a.Login(body.Login, body.Password)
	if err != nil {
		if errors.Is(err, e.HashingFailed) || errors.Is(err, e.FailedToGenerateToken) {
			api.errorResponse(c, http.StatusInternalServerError, err)
		}
		api.errorResponse(c, http.StatusBadRequest, err)
	}

	resp := &response.Login{
		ID:    user.ID,
		Login: user.Login,
		Token: token,
	}

	return c.JSON(http.StatusOK, resp)
}

func (api *V1) Register(c echo.Context) error {
	var body request.Register
	if err := c.Bind(&body); err != nil {
		api.errorResponse(c, http.StatusBadRequest, err)
	}

	if err := api.v.Struct(&body); err != nil {
		api.errorResponse(c, http.StatusBadRequest, err)
	}

	user, err := api.a.Register(body.Login, body.Password)
	if err != nil {
		if errors.Is(err, e.UserAlreadyExists) {
			api.errorResponse(c, http.StatusConflict, err)
		}
		api.errorResponse(c, http.StatusInternalServerError, err)
	}

	resp := &response.Register{
		ID:    user.ID,
		Login: user.Login,
	}

	return c.JSON(http.StatusCreated, resp)
}
