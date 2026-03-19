package viewmodel

import "erp/model"

type FichaTecForm struct {
	FichaID   int
	ModeloID  int
	TecidoRef string
	Tecido    string
	Cor       string
	Cliente   string
	TipoArte  string
	Descricao string
	Custos    model.FichaCustos
	Modelo    model.Modelo
}

func FromFicha(f model.FichaTecnica) FichaTecForm {
	return FichaTecForm{
		FichaID:   f.FichaID,
		ModeloID:  f.ModeloID,
		TecidoRef: f.TecidoRef,
		Tecido:    f.Tecido,
		Cor:       f.Cor,
		Cliente:   f.Cliente,
		TipoArte:  f.TipoArte,
		Descricao: f.Descricao,
		Custos:    f.Custos,
	}
}

func FromFichaComModelo(f model.FichaTecnica, m model.Modelo) FichaTecForm {
	vm := FromFicha(f)
	vm.Modelo = m
	return vm
}
