package service

import (
	"erp/model"
	repository "erp/repository/produtos"
	"errors"
)

// --- CREATE ---
func CriarFicha(f model.FichaTecnica) error {

	if f.ModeloID == 0 {
		return errors.New("modelo é obrigatório")
	}

	if f.Tecido == "" || f.Cor == "" || f.Cliente == "" {
		return errors.New("campos obrigatórios não preenchidos")
	}

	// Calcular total no service
	f.Custos.Total =
		f.Custos.CustoModelo +
			f.Custos.CustoTecido +
			f.Custos.CustoArte

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

func BuscarFichaPorID(id int) (model.FichaTecnica, error) {

	if id == 0 {
		return model.FichaTecnica{}, errors.New("ficha inválida")
	}

	ficha, err := repository.BuscarFichaPorID(id)
	if err != nil {
		return model.FichaTecnica{}, err
	}

	// Regras Futuras

	return ficha, nil
}

// --- UPDATE ---
func AtualizarFichatec(f model.FichaTecnica) error {

	if f.FichaID == 0 {
		return errors.New("ficha inválida")
	}

	// Regra importante:
	f.Custos.Total =
		f.Custos.CustoModelo +
			f.Custos.CustoTecido +
			f.Custos.CustoArte

	return repository.AtualizarFichatec(f)
}

// --- DELETE ---
func DeletarFichatec(id int) error {

	if id == 0 {
		return errors.New("ficha inválida")
	}

	return repository.DeletarFichatec(id)
}
