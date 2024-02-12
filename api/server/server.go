package server

import (
	"github.com/arthurqueiroz4/rinha-de-backend/api/handler"
	"github.com/arthurqueiroz4/rinha-de-backend/repo"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"os"
)

func NewServer(cr *repo.ClienteRepository, tr *repo.TransacaoRepository) {
	port := os.Getenv("APP_PORT")

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("cr", cr)
			c.Set("tr", tr)
			return next(c)
		}
	})
	e.POST("/clientes/:id/transacoes", handler.CriarTransacao)
	e.GET("/clientes/:id/extrato", handler.PegarExtrato)
	e.Logger.Fatal(e.Start(port))
}
