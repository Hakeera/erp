package produtos

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func ProdutosIndex(c echo.Context) error {

	return c.Render(http.StatusOK, "produtos_index", nil)

}
