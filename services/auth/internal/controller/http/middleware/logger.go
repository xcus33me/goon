package middleware

import (
	"auth/pkg/logger"
	"time"

	"github.com/labstack/echo/v4"
)

func Logger(l logger.Interface) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			req := c.Request()
			res := c.Response()

			duration := time.Since(start)

			var logFunc func(message interface{}, args ...interface{})
			if res.Status >= 500 {
				logFunc = l.Info
			} else if res.Status >= 400 {
				logFunc = l.Warn
			} else {
				logFunc = l.Info
			}

			logFunc("HTTP Request",
				"method", req.Method,
				"uri", req.RequestURI,
				"status", res.Status,
				"duration", duration.String(),
				"size", res.Size,
				"remote_ip", c.RealIP(),
				"user_agent", req.UserAgent(),
			)

			return err
		}
	}
}
