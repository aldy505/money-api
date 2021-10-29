package handlers

import (
	"log"
	"money-api/logic/out"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	c.JSON(http.StatusInternalServerError, out.Err{Error: err.Error()})
	log.Println(err)
}
