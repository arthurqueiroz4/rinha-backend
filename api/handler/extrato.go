package handler

import (
	"github.com/arthurqueiroz4/rinha-de-backend/repo"
	"github.com/arthurqueiroz4/rinha-de-backend/types"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"time"
)

func PegarExtrato(c echo.Context) error {
	tr := c.Get("tr").(*repo.TransacaoRepository)
	cr := c.Get("cr").(*repo.ClienteRepository)
	id := c.Param("id")
	if idInt, _ := strconv.Atoi(id); idInt > 5 || idInt < 1 {
		return c.JSON(404, nil)
	}

	transacoes, err := tr.GetLastTransacoes(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "deu erro no GetLastTransacoes")
	}
	cliente, err := cr.GetCliente(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "deu erro no GetCliente")
	}

	extrato := types.Extrato{
		Saldo: types.Saldo{
			Total:       cliente.Saldo,
			DataExtrato: time.Now(),
			Limite:      cliente.Limite,
		},
		UltimasTransacoes: transacoes,
	}

	return c.JSON(http.StatusOK, extrato)
}
