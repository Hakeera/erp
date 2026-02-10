package model

type Produto struct {
	ID      int
	FichaID int

	Tamanho  string
	Linha    string
	Situacao string

	Descricao string
}

type ProdutoComNome struct {
	ID   int
	Nome string
}
