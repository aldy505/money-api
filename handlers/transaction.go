package handlers

import (
	"money-api/logic/account"
	"money-api/logic/auth"
	"money-api/logic/out"
	"money-api/logic/transaction"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func (d *Dependency) GetAllTransaction(c echo.Context) error {
	user := c.Get("UserData").(auth.User)
	t, err := transaction.GetTransactionBy(
		"user",
		transaction.Transaction{
			Sender: account.Account{
				ID: user.ID,
			},
		},
		d.DB,
		c.Request().Context(),
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, t)
}

func (d *Dependency) GetTransactionByFriend(c echo.Context) error {
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

	t, err := transaction.GetTransactionBy(
		"friend",
		transaction.Transaction{
			Sender: account.Account{
				ID: user.ID,
			},
			Recipient: with,
		},
		d.DB,
		c.Request().Context(),
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, t)
}

func (d *Dependency) GetTransactionByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	// Check if transaction exists
	exists, err := transaction.IsTransactionExists(id, d.DB, c.Request().Context(), d.Memory)
	if err != nil {
		return err
	}

	if !exists {
		return c.JSON(http.StatusNotFound, out.Err{Error: "transaction does not found"})
	}

	t, err := transaction.GetTransactionBy("id", transaction.Transaction{ID: id}, d.DB, c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, t)
}

func (d *Dependency) SendMoney(c echo.Context) error {
	var intermediate transaction.Intermediate
	err := c.Bind(intermediate)
	if err != nil {
		return err
	}

	id, err := transaction.CreateTransaction(intermediate, d.DB, c.Request().Context())
	if err != nil {
		return err
	}

	err = transaction.MoveFund(
		intermediate.Sender,
		intermediate.Recipient,
		intermediate.Amount,
		d.DB,
		c.Request().Context(),
	)
	if err != nil {
		return err
	}

	_, err = transaction.UpdateStatus(
		transaction.Transaction{
			ID:     id,
			Status: transaction.StatusSuccess,
		},
		d.DB,
		c.Request().Context(),
	)
	if err != nil {
		return err
	}

	t, err := transaction.GetTransactionBy("id", transaction.Transaction{ID: id}, d.DB, c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, t)
}

func (d *Dependency) RequestMoney(c echo.Context) error {
	var intermediate transaction.Intermediate
	err := c.Bind(intermediate)
	if err != nil {
		return err
	}

	id, err := transaction.CreateRequest(intermediate, d.DB, c.Request().Context())
	if err != nil {
		return err
	}

	t, err := transaction.GetTransactionBy("id", transaction.Transaction{ID: id}, d.DB, c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, t)
}

func (d *Dependency) UpdateStatus(c echo.Context) error {
	var trx transaction.Transaction
	err := c.Bind(trx)
	if err != nil {
		return err
	}

	p := strings.Split(c.Path(), "/")

	switch p[2] {
	case "cancel":
		id, err := transaction.UpdateStatus(transaction.Transaction{
			ID:     trx.ID,
			Status: transaction.StatusCancelled,
		}, d.DB, c.Request().Context())
		if err != nil {
			return err
		}

		t, err := transaction.GetTransactionBy("id", transaction.Transaction{ID: id}, d.DB, c.Request().Context())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, t)
	case "reject":
		id, err := transaction.UpdateStatus(
			transaction.Transaction{
				ID:     trx.ID,
				Status: transaction.StatusCancelled,
			},
			d.DB,
			c.Request().Context(),
		)
		if err != nil {
			return err
		}

		t, err := transaction.GetTransactionBy("id", transaction.Transaction{ID: id}, d.DB, c.Request().Context())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, t)
	default:
		return c.JSON(http.StatusNotFound, out.Err{Error: "not implemented"})
	}
}
