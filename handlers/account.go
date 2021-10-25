package handlers

import (
	"money-api/logic/account"
	"money-api/logic/auth"
	"money-api/logic/out"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (d *Dependency) GetMyAccount(c echo.Context) error {
	user := c.Get("UserData").(auth.User)
	acc, err := account.GetAccountBy("id", account.Account{ID: user.ID}, d.DB, c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, acc)
}

func (d *Dependency) GetFriends(c echo.Context) error {
	user := c.Get("UserData").(auth.User)
	acc, err := account.GetAccountBy("id", account.Account{ID: user.ID}, d.DB, c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, acc.Friends)
}

func (d *Dependency) AddFriend(c echo.Context) error {
	user := c.Get("UserData").(auth.User)
	tag := c.Param("tag")

	// Check if the tag exists
	exists, err := account.IsTagExists(tag, d.DB, c.Request().Context(), d.Memory)
	if err != nil {
		return err
	}

	if !exists {
		return c.JSON(http.StatusNotFound, out.Err{Error: "tag provided does not exists"})
	}

	// Fetch the other guy
	with, err := account.GetAccountBy("tag", account.Account{Tag: tag}, d.DB, c.Request().Context())
	if err != nil {
		return err
	}

	// Check if already friends
	already, err := account.IsFriend(user.ID, with.ID, d.DB, c.Request().Context())
	if err != nil {
		return err
	}

	if already {
		return c.JSON(http.StatusBadRequest, out.Msg{Message: "already friends"})
	}

	acc, err := account.AddFriend(user.ID, with.ID, d.DB, c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, acc)
}

func (d *Dependency) RemoveFriend(c echo.Context) error {
	return nil
}
