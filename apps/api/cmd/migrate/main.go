package main

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"log"
	"os"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL requerida")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("configuración de migraciones inválida")
	}
	defer db.Close()
	if err = goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}
	ruta := os.Getenv("MIGRATIONS_DIR")
	if ruta == "" {
		ruta = "/app/db/migrations"
	}
	if err = goose.Up(db, ruta); err != nil {
		log.Fatal("migración fallida: ", err)
	}
}
