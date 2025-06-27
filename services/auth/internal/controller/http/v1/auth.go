package v1

import (
	"auth/internal/controller/http/v1/request"
	"auth/internal/controller/http/v1/response"
	e "auth/pkg/errors"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (api *V1) Login(c echo.Context) error {
	var body request.Login
	if err := c.Bind(&body); err != nil {
		return api.errorResponse(c, http.StatusBadRequest, err)
	}

	if err := api.v.Struct(&body); err != nil {
		return api.errorResponse(c, http.StatusBadRequest, err)
	}

	user, token, err := api.a.Login(c.Request().Context(), body.Login, body.Password)
	if err != nil {
		if errors.Is(err, e.HashingFailed) || errors.Is(err, e.FailedToGenerateToken) {
			return api.errorResponse(c, http.StatusInternalServerError, err)
		}
		return api.errorResponse(c, http.StatusBadRequest, err)
	}

	cookie := new(http.Cookie)
	cookie.Name = "jwt_token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(time.Hour * 24)
	cookie.Path = "/"
	cookie.HttpOnly = true

	c.SetCookie(cookie)

	resp := &response.Login{
		ID:    user.ID,
		Login: user.Login,
	}

	return c.JSON(http.StatusOK, resp)
}

func (api *V1) Register(c echo.Context) error {
	var body request.Register
	if err := c.Bind(&body); err != nil {
		return api.errorResponse(c, http.StatusBadRequest, err)
	}

	if err := api.v.Struct(&body); err != nil {
		return api.errorResponse(c, http.StatusBadRequest, err)
	}

	user, err := api.a.Register(c.Request().Context(), body.Login, body.Password)
	if err != nil {
		if errors.Is(err, e.UserAlreadyExists) {
			return api.errorResponse(c, http.StatusConflict, err)
		}
		return api.errorResponse(c, http.StatusInternalServerError, err)
	}

	resp := &response.Register{
		ID:    user.ID,
		Login: user.Login,
	}

	return c.JSON(http.StatusCreated, resp)
}
