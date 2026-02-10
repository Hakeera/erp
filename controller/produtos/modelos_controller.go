package produtos

import (
	"erp/model"
	repository "erp/repository/produtos"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func ModelosIndex(c echo.Context) error {

	return c.Render(http.StatusOK, "modelos_index", nil)

}

func ModelosForm(c echo.Context) error {
	idStr := c.QueryParam("id")

	if idStr == "" {
		return c.Render(http.StatusOK, "modelos_form", model.Modelo{})
	}

	id, _ := strconv.Atoi(idStr)

	modelo, err := repository.BuscarModelo(id)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "modelos_form", modelo)
}

func ModelosList(c echo.Context) error {

	modelos, err := repository.ListarModelos()
	if err != nil {
		return err
	}

	data := struct {
		Modelos []model.Modelo
	}{
		Modelos: modelos,
	}

	return c.Render(http.StatusOK, "modelos_list", data)
}

func ModelosCreate(c echo.Context) error {

	modelo := model.Modelo{
		Nome:      c.FormValue("nome"),
		Linha:     c.FormValue("linha"),
		Descricao: c.FormValue("descricao"),
	}

	modelo.Corte, _ = strconv.Atoi(c.FormValue("corte"))
	modelo.Costura, _ = strconv.Atoi(c.FormValue("costura"))
	modelo.Acabamento, _ = strconv.Atoi(c.FormValue("acabamento"))
	modelo.Aviamento, _ = strconv.Atoi(c.FormValue("aviamento"))

	if err := repository.CriarModelo(modelo); err != nil {
		return err
	}

	// recarrega só a lista
	return ModelosList(c)
}
