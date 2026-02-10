package produtos

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func ProdutosCatalogo(c echo.Context) error {

	return c.Render(http.StatusOK, "produtos_catalogo", nil)

}
