package produtos

import (
	"erp/model"
	"erp/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// --- Página Inicial com todas as partes carregadas e renderizadas ---
func FichatecIndex(c echo.Context) error {

	fichas, err := service.ListarFichatec()
	if err != nil {
		return err
	}

	data := struct {
		Fichas   []model.FichaTecnica
		FichaTec model.FichaTecnica
	}{
		Fichas: fichas,
		FichaTec: model.FichaTecnica{
			Custos: model.FichaCustos{},
		},
	}

	return c.Render(http.StatusOK, "fichatec_index", data)
}

// --- CREATE ---
func CriarFicha(c echo.Context) error {

	f := model.FichaTecnica{
		TecidoRef: c.FormValue("tecido_ref"),
		Tecido:    c.FormValue("tecido"),
		Cor:       c.FormValue("cor"),
		Cliente:   c.FormValue("cliente"),
		Descricao: c.FormValue("descricao"),
		TipoArte:  c.FormValue("tipo_arte"),
	}

	f.ModeloID, _ = strconv.Atoi(c.FormValue("modelo_id"))

	f.Custos.CustoModelo, _ = strconv.Atoi(c.FormValue("custo_modelo"))
	f.Custos.CustoTecido, _ = strconv.Atoi(c.FormValue("custo_tecido"))
	f.Custos.CustoArte, _ = strconv.Atoi(c.FormValue("custo_arte"))

	if err := service.CriarFicha(f); err != nil {
		return c.String(400, err.Error())
	}

	return FichatecIndex(c)
}

// --- READ ---
func FichatecForm(c echo.Context) error {

	idStr := c.QueryParam("id")
	// Se não tiver ID → formulário vazio (novo)
	if idStr == "" {
		return c.Render(http.StatusOK, "fichatec_form", model.FichaTecnica{
			Custos: model.FichaCustos{},
		})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(400, "ID inválido")
	}

	ficha, err := service.BuscarFichaPorID(id)
	if err != nil {
		return c.String(404, err.Error())
	}

	return c.Render(http.StatusOK, "fichatec_form", ficha)
}

func ListarFichatec(c echo.Context) error {

	return nil
}

// --- UPDATE ---
func AtualizarFichatec(c echo.Context) error {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(400, "ID inválido")
	}

	ficha := model.FichaTecnica{
		FichaID:   id,
		TecidoRef: c.FormValue("tecido_ref"),
		Tecido:    c.FormValue("tecido"),
		Cor:       c.FormValue("cor"),
		Cliente:   c.FormValue("cliente"),
		Descricao: c.FormValue("descricao"),
		TipoArte:  c.FormValue("tipo_arte"),
	}

	ficha.ModeloID, _ = strconv.Atoi(c.FormValue("modelo_id"))

	ficha.Custos = model.FichaCustos{}
	ficha.Custos.CustoModelo, _ = strconv.Atoi(c.FormValue("custo_modelo"))
	ficha.Custos.CustoTecido, _ = strconv.Atoi(c.FormValue("custo_tecido"))
	ficha.Custos.CustoArte, _ = strconv.Atoi(c.FormValue("custo_arte"))

	if err := service.AtualizarFichatec(ficha); err != nil {
		return c.String(400, err.Error())
	}

	return FichatecIndex(c)
}

// --- DELETE ---
func DeletarFichatec(c echo.Context) error {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(400, "ID inválido")
	}

	if err := service.DeletarFichatec(id); err != nil {
		return c.String(400, err.Error())
	}

	return FichatecIndex(c)
}
