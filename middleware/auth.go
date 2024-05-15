package middleware

import (
	"errors"
	"strings"

	"halo_suster/crypto"
	"halo_suster/responses"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func ItAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get(echo.HeaderAuthorization)
		if auth == "" {
			return responses.NewUnauthorizedError("token not found")
		}

		splitted := strings.Split(auth, " ")
		if len(splitted) != 2 || splitted[0] != "Bearer" {
			return responses.NewUnauthorizedError("invalid token")
		}

		payload, err := crypto.VerifyToken(splitted[1])
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				return responses.NewUnauthorizedError("token expired")
			}
			return responses.NewUnauthorizedError(err.Error())
		}

		if !strings.HasPrefix(payload.Nip, "615") {
			return responses.NewUnauthorizedError("user is not an it staff")
		}

		c.Set("user_id", payload.Id)
		c.Set("nip", payload.Nip)
		c.Set("name", payload.Name)

		return next(c)
	}
}

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get(echo.HeaderAuthorization)
		if auth == "" {
			return responses.NewUnauthorizedError("token not found")
		}

		splitted := strings.Split(auth, " ")
		if len(splitted) != 2 || splitted[0] != "Bearer" {
			return responses.NewUnauthorizedError("invalid token")
		}

		payload, err := crypto.VerifyToken(splitted[1])
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				return responses.NewUnauthorizedError("token expired")
			}
			return responses.NewUnauthorizedError(err.Error())
		}

		c.Set("user_id", payload.Id)
		c.Set("nip", payload.Nip)
		c.Set("name", payload.Name)

		return next(c)
	}
}
