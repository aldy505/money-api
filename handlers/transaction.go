package handlers

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func (d *Dependency) GetAllTransaction(c echo.Context) error {
	return nil
}

func (d *Dependency) GetTransactionByFriend(c echo.Context) error {
	return nil
}

func (d *Dependency) GetTransactionByID(c echo.Context) error {
	return nil
}

func (d *Dependency) SendMoney(c echo.Context) error {
	return nil
}

func (d *Dependency) RequestMoney(c echo.Context) error {
	return nil
}

func (d *Dependency) UpdateStatus(c echo.Context) error {
	p := strings.Split(c.Path(), "/")

	switch p[2] {
	case "cancel":
		return nil
	case "reject":
		return nil
	default:
		return nil
	}
}
