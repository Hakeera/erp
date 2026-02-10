package model

type Modelo struct {
	ID    int
	Nome  string
	Linha string

	Corte      int
	Costura    int
	Acabamento int
	Aviamento  int

	ConsumoPorGrade map[string]float64
	Descricao       string
}
