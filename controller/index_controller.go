package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func IndexPage(c echo.Context) error {

	return c.Render(http.StatusOK, "index", nil)
}
