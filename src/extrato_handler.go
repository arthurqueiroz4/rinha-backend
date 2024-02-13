package main

import (
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"time"
)

func PegarExtrato(c echo.Context) error {
	id := c.Param("id")
	if idInt, _ := strconv.Atoi(id); idInt > 5 || idInt < 1 {
		return c.JSON(404, nil)
	}

	transacoes, err := PegarUltimasTransacoesDB(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "deu erro no GetLastTransacoes")
	}
	cliente, err := PegarClienteDB(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "deu erro no PegarClienteDB")
	}

	extrato := Extrato{
		Saldo: Saldo{
			Total:       cliente.Saldo,
			DataExtrato: time.Now(),
			Limite:      cliente.Limite,
		},
		UltimasTransacoes: transacoes,
	}

	return c.JSON(http.StatusOK, extrato)
}
