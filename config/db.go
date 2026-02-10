package config

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DB   *pgxpool.Pool
	once sync.Once
)

func InitDB() {
	once.Do(func() {
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			log.Fatal("DATABASE_URL não definida")
		}

		var err error
		DB, err = pgxpool.New(context.Background(), dsn)
		if err != nil {
			log.Fatalf("Erro ao criar pool: %v", err)
		}

		if err := DB.Ping(context.Background()); err != nil {
			log.Fatalf("Erro ao conectar no banco: %v", err)
		}

		log.Println("PostgreSQL conectado com sucesso")
	})
}

func GetDB() *pgxpool.Pool {
	if DB == nil {
		log.Fatal("Banco não inicializado")
	}
	return DB
}
