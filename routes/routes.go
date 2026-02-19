package routes

import (
	"erp/controller"
	"erp/controller/produtos"

	"github.com/labstack/echo/v4"
)

func SetUpRoutes(e *echo.Echo) {

	// Rotas das Páginas Iniciais
	e.GET("/", controller.IndexPage)

	// Página Inicial de Produtos
	e.GET("/produtos/index", produtos.ProdutosIndex)
	e.GET("/produtos/catalogo", produtos.ProdutosCatalogo)

	// Modelos
	e.GET("/produtos/modelos", produtos.ModelosIndex)         // Carrega Modelos (Lista e Form)
	e.POST("/produtos/modelos", produtos.CriarModelo)         // Criar Modelo
	e.GET("/produtos/modelos/form", produtos.ModelosForm)     // Formulário de Edição de Modelos
	e.PUT("/produtos/modelos/:id", produtos.AtualizarModelo)  // Atualizar Modelo
	e.DELETE("/produtos/modelos/:id", produtos.DeletarModelo) // Excluir Modelo

	// Ficha Técnica
	e.GET("/produtos/fichatec", produtos.FichatecIndex)          // Carrega Fichatec (Lista e Form)
	e.GET("/produtos/fichatec/form", produtos.FichatecForm)      // Formulário de Edição de Modelos
	e.PUT("/produtos/fichatec/:id", produtos.AtualizarFichatec)  // Atualizar Fichatec
	e.DELETE("/produtos/fichatec/:id", produtos.DeletarFichatec) // Excluir Fichatec

	// TODO: Criação de Fichas deve obter o Nome de Modelo e custos do modelo automaticamente
	// Criar Ficha deve trazer uma prévisualização dos produtos que vão ser criados
	// Ficha_Tabela deve ser implementada com MarkUp e tabelas de preços para isso
	e.POST("/produtos/fichatec", produtos.CriarFicha) // Cria FichaTec

	// FIXME: Essas rotas são necessárias sendo que Index já faz a listagem?
	// e.GET("/produtos/modelos/list", produtos.ListarModelos)   // Listagem de Modelos
	// e.GET("/produtos/fichatec/list", produtos.ListarFichatec) // Listagem de Fichatec

	// e.GET("/clientes"", controller.)

	// e.GET("/pedidos", controller.)

	// e.GET("/producao", controller.)

}
