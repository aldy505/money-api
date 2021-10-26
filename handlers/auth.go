package handlers

import (
	"money-api/logic/auth"
	"money-api/logic/out"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (d *Dependency) Login(c echo.Context) error {
	var user auth.User
	err := c.Bind(&user)
	if err != nil {
		return err
	}

	check, err := auth.CheckIfUserExists(user, d.DB, c.Request().Context())
	if err != nil {
		return err
	}

	if !check {
		return c.JSON(http.StatusBadRequest, out.Err{Error: "email does not exists"})
	}

	fetch, err := auth.GetUserByEmail(user.Email, d.DB, c.Request().Context())
	if err != nil {
		return err
	}

	validate, err := auth.VerifyPassword(user.Password, fetch.Password)
	if err != nil {
		return err
	}

	if !validate {
		return c.JSON(http.StatusForbidden, out.Err{Error: "password do not match"})
	}

	token, err := auth.GenerateJWT(d.JWTSecret, fetch)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, out.AuthToken{Token: token})
}

func (d *Dependency) Signup(c echo.Context) error {
	var user auth.User
	err := c.Bind(&user)
	if err != nil {
		return err
	}

	check, err := auth.CheckIfUserExists(user, d.DB, c.Request().Context())
	if err != nil {
		return err
	}

	if check {
		return c.JSON(http.StatusBadRequest, out.Err{Error: "user already exists"})
	}

	pass, err := auth.GeneratePassword(user.Password)
	if err != nil {
		return err
	}

	finalUser := auth.User{
		Password: pass,
		Name:     user.Name,
		Email:    user.Email,
		Address:  user.Address,
	}
	u, err := auth.RegisterUser(finalUser, d.DB, c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, auth.User{ID: u.ID, Name: u.Name, Email: u.Email})
}

func (d *Dependency) Verify(c echo.Context) error {
	// Get token from headers
	bearer := strings.Replace(c.Request().Header.Get("authorization"), "Bearer ", "", 1)
	validate, err := auth.VerifyJWT(d.JWTSecret, bearer)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, validate)
}
