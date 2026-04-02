package viewmodel

import "erp/model"

type FichaTecForm struct {
	FichaID     int
	ModeloID    int
	TecidoRef   string
	Tecido      string
	Cor         string
	Cliente     string
	TipoArte    string
	Descricao   string
	CustosGrade []model.FichaCustoGrade
	Custos      model.FichaCustos
	Modelo      model.Modelo
}

func FromFicha(f model.FichaTecnica) FichaTecForm {
	return FichaTecForm{
		FichaID:     f.FichaID,
		ModeloID:    f.ModeloID,
		TecidoRef:   f.TecidoRef,
		Tecido:      f.Tecido,
		Cor:         f.Cor,
		Cliente:     f.Cliente,
		TipoArte:    f.TipoArte,
		CustosGrade: f.CustosGrade,
		Descricao:   f.Descricao,
		Custos:      f.Custos,
	}
}

func FromFichaComModelo(f model.FichaTecnica, m model.Modelo) FichaTecForm {
	vm := FromFicha(f)
	vm.Modelo = m
	return vm
}

func FromModeloParaFichaForm(m model.Modelo) FichaTecForm {

	custoModelo := m.Corte + m.Costura + m.Acabamento + m.Aviamento

	return FichaTecForm{
		ModeloID: m.ID,
		Modelo:   m,
		Custos: model.FichaCustos{
			CustoModelo: custoModelo,
		},
	}
}

// Obter Custos de Tecido de cada grade
func (vm FichaTecForm) GetCustoGrade(nome string) int {
	for _, g := range vm.CustosGrade {
		if g.Grade == nome {
			return g.Custo
		}
	}
	return 0
}
