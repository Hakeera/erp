package produtos

import (
	"erp/model"
	"erp/service"
	"erp/viewmodel"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// --- Página Inicial com todas as partes carregadas e renderizadas ---
// FichatecIndex — passa viewmodel no FichaTec
func FichatecIndex(c echo.Context) error {
	fichas, err := service.ListarFichatec()
	if err != nil {
		return err
	}

	// monta lista com modelo resolvido
	fichasForms := make([]viewmodel.FichaTecForm, 0, len(fichas))
	for _, f := range fichas {
		vm := viewmodel.FromFicha(f)
		if f.ModeloID != 0 {
			modelo, err := service.BuscarModeloPorID(f.ModeloID)
			if err == nil {
				vm.Modelo = modelo
			}
		}
		fichasForms = append(fichasForms, vm)
	}

	data := struct {
		Fichas   []viewmodel.FichaTecForm
		FichaTec viewmodel.FichaTecForm
	}{
		Fichas:   fichasForms,
		FichaTec: viewmodel.FichaTecForm{},
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
		log.Println("ERRO REPO CRIAR FICHA:", err)
		return c.String(400, err.Error())
	}

	return FichatecIndex(c)
}

// --- READ ---

// FichatecForm — vazio ou edição
func FichatecForm(c echo.Context) error {
	idStr := c.QueryParam("id")
	if idStr == "" {
		return c.Render(http.StatusOK, "fichatec_form", viewmodel.FichaTecForm{})
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(400, "ID inválido")
	}
	ficha, err := service.BuscarFichaPorID(id)
	if err != nil {
		return c.String(404, err.Error())
	}
	return c.Render(http.StatusOK, "fichatec_form", viewmodel.FromFicha(ficha))
}

// ModeloPraFichaTec — selecionou modelo no overlay
// Seleciona o Modelo para o Formulário da Fichatec
func ModeloPraFichaTec(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(400, "ID inválido")
	}

	modelo, err := service.BuscarModeloPorID(id)
	if err != nil {
		return c.String(500, err.Error())
	}

	vm := viewmodel.FromModeloParaFichaForm(modelo)

	return c.Render(http.StatusOK, "fichatec_form", vm)
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

// Através da requisição filtra os modelos de acordo com Linha e Nome
// Retorna listagem de modelos para o render no overlay
func FiltrarModelos(c echo.Context) error {

	nome := c.QueryParam("nome")
	linha := c.QueryParam("linha")

	modelos, err := service.ModelosComFiltro(nome, linha)
	if err != nil {
		return c.String(500, err.Error())
	}

	data := struct {
		Modelos []model.Modelo
	}{
		Modelos: modelos,
	}

	return c.Render(http.StatusOK, "fichatec_modelos_search", data)
}
