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
	e.GET("/produtos/modelos", produtos.ModelosIndex)         // Setor de Modelos
	e.GET("/produtos/modelos/list", produtos.ListarModelos)   // Listagem de Modelos
	e.GET("/produtos/modelos/form", produtos.ModelosForm)     // Formulário de Edição de Modelos
	e.POST("/produtos/modelos", produtos.CriarModelo)         // Criar Modelo
	e.PUT("/produtos/modelos/:id", produtos.AtualizarModelo)  // Atualizar Modelo
	e.DELETE("/produtos/modelos/:id", produtos.DeletarModelo) // Excluir Modelo

	// Ficha Técnica
	e.GET("/produtos/fichatec", produtos.FichaTecIndex)

	// e.GET("/clientes"", controller.)

	// e.GET("/pedidos", controller.)

	// e.GET("/producao", controller.)

}
