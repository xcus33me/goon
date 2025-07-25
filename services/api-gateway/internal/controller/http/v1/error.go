package v1

import "github.com/labstack/echo/v4"

func (api *V1) errorResponse(c echo.Context, code int, err error) error {
	api.l.Error(err)
	return c.JSON(code, map[string]string{"error": err.Error()})
}
