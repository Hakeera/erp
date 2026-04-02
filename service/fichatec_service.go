package service

import (
	"erp/model"
	repository "erp/repository/produtos"
	"errors"
	"fmt"
)

// --- CREATE ---
func CriarFicha(f model.FichaTecnica) error {

	if f.ModeloID == 0 {
		return errors.New("modelo é obrigatório")
	}

	if f.Tecido == "" || f.Cor == "" || f.Cliente == "" {
		return errors.New("campos obrigatórios não preenchidos")
	}

	if len(f.CustosGrade) == 0 {
		return errors.New("custos por grade são obrigatórios")
	}

	fmt.Println("CUSTO GRADE:", f.CustosGrade)
	// valida se todas as grades têm valor
	for _, cg := range f.CustosGrade {
		if cg.Custo <= 0 {
			return fmt.Errorf("custo inválido para grade %s", cg.Grade)
		}
	}

	return repository.CriarFicha(f)
}

// --- READ ---
func ListarFichatec() ([]model.FichaTecnica, error) {

	fichas, err := repository.ListarFichas()
	if err != nil {
		return nil, err
	}

	// Regras Futuras
	return fichas, nil
}

func BuscarFichaPorID(id int) (model.FichaTecnica, model.Modelo, error) {
	if id == 0 {
		return model.FichaTecnica{}, model.Modelo{}, errors.New("ficha inválida")
	}

	ficha, modelo, err := repository.BuscarFichaPorID(id)
	if err != nil {
		return model.FichaTecnica{}, model.Modelo{}, err
	}

	return ficha, modelo, nil
}

// --- UPDATE ---
func AtualizarFichatec(f model.FichaTecnica) error {

	if f.FichaID == 0 {
		return errors.New("ficha inválida")
	}

	return repository.AtualizarFichatec(f)
}

// --- DELETE ---
func DeletarFichatec(id int) error {

	if id == 0 {
		return errors.New("ficha inválida")
	}

	return repository.DeletarFichatec(id)
}
