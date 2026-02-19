package produtos

import (
	"erp/model"
	"erp/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// --- Página Inicial com todas as partes carregadas e renderizadas ---
func ModelosIndex(c echo.Context) error {

	modelos, err := service.ListarModelos()
	if err != nil {
		return err
	}

	data := struct {
		Modelos []model.Modelo
		Modelo  model.Modelo
	}{
		Modelos: modelos,
		Modelo:  model.Modelo{}, // form vazio
	}

	return c.Render(http.StatusOK, "modelos_index", data)
}

// --- CREATE ---
func CriarModelo(c echo.Context) error {

	modelo := model.Modelo{
		Nome:      c.FormValue("nome"),
		Linha:     c.FormValue("linha"),
		Descricao: c.FormValue("descricao"),
	}

	modelo.Corte, _ = strconv.Atoi(c.FormValue("corte"))
	modelo.Costura, _ = strconv.Atoi(c.FormValue("costura"))
	modelo.Acabamento, _ = strconv.Atoi(c.FormValue("acabamento"))
	modelo.Aviamento, _ = strconv.Atoi(c.FormValue("aviamento"))

	modelo.ConsumoPorGrade = make(map[string]float64)

	categorias := map[string]string{
		"INFANTIL": "grade_infantil",
		"JUVENIL":  "grade_juvenil",
		"ADULTO":   "grade_adulto",
		"EXTRA":    "grade_extra",
	}

	for chaveJSON, campoForm := range categorias {

		valorStr := c.FormValue(campoForm)
		if valorStr == "" {
			continue
		}

		valor, err := strconv.ParseFloat(valorStr, 64)
		if err != nil {
			return err
		}

		modelo.ConsumoPorGrade[chaveJSON] = valor
	}

	if err := service.CriarModelo(modelo); err != nil {
		return c.String(400, err.Error())
	}
	return ModelosIndex(c)
}

// --- READ ---
func ListarModelos(c echo.Context) error {

	modelos, err := service.ListarModelos()
	if err != nil {
		return err
	}

	data := struct {
		Modelos []model.Modelo
	}{
		Modelos: modelos,
	}

	return c.Render(http.StatusOK, "modelos_index", data)
}

func ModelosForm(c echo.Context) error {

	idStr := c.QueryParam("id")
	// Se não tiver ID → formulário vazio (novo)
	if idStr == "" {
		return c.Render(http.StatusOK, "modelos_form", model.Modelo{})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(400, "ID inválido")
	}

	modelo, err := service.BuscarModeloPorID(id)
	if err != nil {
		return c.String(404, err.Error())
	}

	return c.Render(http.StatusOK, "modelos_form", modelo)
}

// --- UPDATE---
func AtualizarModelo(c echo.Context) error {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(400, "ID inválido")
	}

	modelo := model.Modelo{
		ID:        id,
		Nome:      c.FormValue("nome"),
		Linha:     c.FormValue("linha"),
		Descricao: c.FormValue("descricao"),
	}

	modelo.Corte, _ = strconv.Atoi(c.FormValue("corte"))
	modelo.Costura, _ = strconv.Atoi(c.FormValue("costura"))
	modelo.Acabamento, _ = strconv.Atoi(c.FormValue("acabamento"))
	modelo.Aviamento, _ = strconv.Atoi(c.FormValue("aviamento"))

	// Consumo por categoria
	modelo.ConsumoPorGrade = make(map[string]float64)

	categorias := map[string]string{
		"INFANTIL": "grade_infantil",
		"JUVENIL":  "grade_juvenil",
		"ADULTO":   "grade_adulto",
		"EXTRA":    "grade_extra",
	}

	for chaveJSON, campoForm := range categorias {

		valorStr := c.FormValue(campoForm)
		if valorStr == "" {
			continue
		}

		valor, err := strconv.ParseFloat(valorStr, 64)
		if err != nil {
			return err
		}

		modelo.ConsumoPorGrade[chaveJSON] = valor
	}

	if err := service.AtualizarModelo(modelo); err != nil {
		return c.String(400, err.Error())
	}

	// Atualiza apenas a lista
	return ModelosIndex(c)
}

// --- DELETE ---
func DeletarModelo(c echo.Context) error {

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(400, "ID inválido")
	}

	if err := service.DeletarModelo(id); err != nil {
		return c.String(400, err.Error())
	}

	return ModelosIndex(c)
}
