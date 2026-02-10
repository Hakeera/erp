package routes

import (
	"erp/controller"

	"github.com/labstack/echo/v4"
)

func SetUpRoutes(e *echo.Echo) {

	// Rotas das Páginas Iniciais
	e.GET("/", controller.IndexPage)

	// e.GET("/produtos", controller.)

	// e.GET("/clientes"", controller.)

	// e.GET("/pedidos", controller.)

	// e.GET("/producao", controller.)

}
