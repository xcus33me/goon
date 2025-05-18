package middleware

import (
	"auth/internal/entity"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userCtx"
	cookieName          = "jwt_token"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int    `json:"id"`
	Login  string `json:"login"`
}

func Auth(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var tokenString string

			cookie, err := c.Cookie(cookieName)
			if err == nil && cookie != nil && cookie.Value != "" {
				tokenString = cookie.Value
			} else {
				authHeader := c.Request().Header.Get(authorizationHeader)
				if authHeader == "" {
					return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token: no cookie or auth header"})
				}

				headerParts := strings.Split(authHeader, " ")
				if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
					return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid auth header format"})
				}
				tokenString = headerParts[1]
			}

			if tokenString == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "token is empty"})
			}

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method: " + token.Header["alg"].(string))
				}
				return []byte(jwtSecret), nil
			})

			if err != nil {
				if errors.Is(err, jwt.ErrTokenExpired) {
					return c.JSON(http.StatusUnauthorized, map[string]string{"error": "token expired"})
				}
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token: " + err.Error()})
			}

			if !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "token is invalid"})
			}

			user := &entity.User{
				ID:    claims.UserID,
				Login: claims.Login,
			}
			c.Set(userCtx, user)

			return next(c)
		}
	}
}

func GetUserFromContext(c echo.Context) (*entity.User, error) {
	userVal := c.Get(userCtx)
	if userVal == nil {
		return nil, errors.New("user not found in context: context key missing")
	}

	user, ok := userVal.(*entity.User)
	if !ok {
		return nil, errors.New("user not found in context: type assertion failed")
	}
	return user, nil
}
