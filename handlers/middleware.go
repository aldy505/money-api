package handlers

import (
	"net/http"
	"strings"

	"money-api/logic/auth"
	"money-api/logic/out"

	"github.com/labstack/echo/v4"
)

func (d *Dependency) MustLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get token from headers
		bearer := strings.Replace(c.Request().Header.Get("authorization"), "Bearer ", "", 1)
		data, err := auth.VerifyJWT(d.JWTSecret, bearer)
		if err != nil {
			return c.JSON(http.StatusForbidden, out.Err{Error: err.Error()})
		}

		c.Set("UserData", data)

		return next(c)
	}
}
