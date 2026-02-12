package service

import (
	"erp/model"
	repository "erp/repository/produtos"
	"errors"
)

var MapeamentoGrade = map[string]string{
	"00": "INFANTIL",
	"02": "INFANTIL",
	"04": "INFANTIL",
	"06": "INFANTIL",
	"08": "INFANTIL",

	"10": "JUVENIL",
	"12": "JUVENIL",
	"14": "JUVENIL",
	"16": "JUVENIL",

	"PP": "ADULTO",
	"P":  "ADULTO",
	"M":  "ADULTO",
	"G":  "ADULTO",
	"GG": "ADULTO",

	"XG":  "EXTRA",
	"XGG": "EXTRA",
}

// --- CREATE ---
func CriarModelo(m model.Modelo) error {

	if m.Nome == "" {
		return errors.New("nome é obrigatório")
	}

	if len(m.ConsumoPorGrade) == 0 {
		return errors.New("é necessário informar pelo menos um consumo de grade")
	}

	return repository.CriarModelo(m)
}

// --- READ ---
func ListarModelos() ([]model.Modelo, error) {

	modelos, err := repository.ListarModelos()
	if err != nil {
		return nil, err
	}

	// Aqui entram regras futuras
	// validar se está ativo
	// aplicar transformação
	// enriquecer dados
	return modelos, nil
}

func BuscarModeloPorID(id int) (model.Modelo, error) {

	if id == 0 {
		return model.Modelo{}, errors.New("modelo inválido")
	}

	modelo, err := repository.BuscarModeloPorID(id)
	if err != nil {
		return model.Modelo{}, err
	}

	// Aqui entram regras futuras
	// validar se está ativo
	// aplicar transformação
	// enriquecer dados

	return modelo, nil
}

// --- UPDATE ---
func AtualizarModelo(m model.Modelo) error {

	if m.ID == 0 {
		return errors.New("modelo inválido")
	}

	if m.Nome == "" {
		return errors.New("nome é obrigatório")
	}

	return repository.AtualizarModelo(m)
}

// --- DELETE ---
func DeletarModelo(id int) error {

	if id == 0 {
		return errors.New("modelo inválido")
	}

	return repository.DeletarModelo(id)
}

// --- Funções Auxiliares ---
func CalcularConsumoTotal(modelo model.Modelo, pedido map[string]int) float64 {

	total := 0.0

	for tamanho, quantidade := range pedido {

		categoria := MapeamentoGrade[tamanho]

		consumoBase, ok := modelo.ConsumoPorGrade[categoria]
		if !ok {
			continue
		}

		total += consumoBase * float64(quantidade)
	}

	return total
}
