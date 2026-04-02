package main

import (
	"context"
	"erp/config"
	"erp/routes"
	"fmt"
	"io"
	"log"
	"text/template"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

// TemplateRenderer para Echo
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders tamplates with data
func (t *TemplateRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	// Debug: verificar se o template existe
	tmpl := t.templates.Lookup(name)
	if tmpl == nil {
		log.Printf("❌ Template '%s' não encontrado!", name)
		// Listar todos os templates disponíveis
		for _, t := range t.templates.Templates() {
			log.Printf("📄 Template disponível: %s", t.Name())
		}
		return fmt.Errorf("template %s não encontrado", name)
	}

	log.Printf("✅ Renderizando template: %s", name)
	log.Printf("📊 Dados: %+v", data)

	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {

	// Carrega variáveis de ambiente do .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env não encontrado, usando variáveis do sistema")
	}

	// Inicializa banco
	config.InitDB()
	log.Println("✅ Banco inicializado")

	// Teste rápido de vida
	ctx := context.Background()
	var now time.Time
	err := config.GetDB().
		QueryRow(ctx, "SELECT now()").
		Scan(&now)
	if err != nil {
		log.Fatalf("❌ Banco não respondeu: %v", err)
	}
	log.Printf("🕒 Banco respondeu em: %s", now)

	// Templates
	templates := template.New("")
	template.Must(templates.ParseGlob("view/*/*.html"))
	template.Must(templates.ParseGlob("view/*/*/*.html"))
	renderer := &TemplateRenderer{
		templates: templates,
	}
	for _, tmpl := range renderer.templates.Templates() {
		log.Printf("📄 Template carregado: %s", tmpl.Name())
	}

	// Echo
	e := echo.New()
	e.Static("/static", "view/static")
	e.Renderer = renderer

	// Rotas
	routes.SetUpRoutes(e)

	// Server
	e.Logger.Fatal(e.Start(":8080"))
}
