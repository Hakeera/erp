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
	e.GET("/produtos/modelos", produtos.ModelosIndex)     // Setor de Modelos
	e.GET("/produtos/modelos/list", produtos.ModelosList) // Listagem de Modelos
	e.GET("/produtos/modelos/form", produtos.ModelosForm) // Formulário de Edição de Modelos
	e.POST("/produtos/modelos", produtos.ModelosCreate)   // Criar Modelo
	e.PUT("/produtos/modelos/:id", produtos.ModelosCreate)

	// Ficha Técnica
	e.GET("/produtos/fichatec", produtos.FichaTecIndex)

	// e.GET("/clientes"", controller.)

	// e.GET("/pedidos", controller.)

	// e.GET("/producao", controller.)

}
