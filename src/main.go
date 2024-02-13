package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"os"
)

var db *pgxpool.Pool
var cache map[string]*Cliente

func main() {
	var err error
	db, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	err = db.Ping(context.Background())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	cache = make(map[string]*Cliente)

	port := os.Getenv("APP_PORT")
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.POST("/clientes/:id/transacoes", CriarTransacao)
	e.GET("/clientes/:id/extrato", PegarExtrato)
	e.Logger.Fatal(e.Start(port))
}
