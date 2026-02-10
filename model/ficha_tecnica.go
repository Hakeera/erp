package model

type FichaTecnica struct {
	ID       int
	ModeloID int

	TecidoRef string
	Tecido    string
	Cor       string
	Cliente   string

	TipoArte  string
	Descricao string

	Custos  FichaCustos
	Tabelas []FichaTabela
}

type FichaCustos struct {
	FichaID int

	CustoModelo int
	CustoTecido int
	CustoArte   int

	CustosExtras map[string]int

	Total int
}

type FichaTabela struct {
	ID      int
	FichaID int

	Nome          string
	Markup        int
	PrecoOverride *int
}
