package produtos

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func FichaTecIndex(c echo.Context) error {

	return c.Render(http.StatusOK, "fichatec_index", nil)

}
