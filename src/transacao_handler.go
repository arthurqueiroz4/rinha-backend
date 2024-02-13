package main

import (
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func CriarTransacao(c echo.Context) error {
	id := c.Param("id")
	if idInt, _ := strconv.Atoi(id); idInt > 5 || idInt < 1 {
		return c.JSON(404, nil)
	}

	transacao := &Transacao{}
	transacao.ClienteID, _ = strconv.Atoi(id)
	if err := c.Bind(transacao); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, nil)
	}

	novoSaldo := transacao.Valor
	if transacao.Tipo != "d" && transacao.Tipo != "c" {
		return c.JSON(http.StatusUnprocessableEntity, nil)
	} else if transacao.Tipo == "d" {
		novoSaldo *= -1
	}

	if descLen := len(transacao.Descricao); descLen < 1 || descLen > 10 {
		return c.JSON(http.StatusUnprocessableEntity, nil)
	}

	cliente, err := PegarClienteDB(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "deu erro no PegarClienteDB")
	}
	if (cliente.Saldo + transacao.Valor) < -cliente.Limite {
		return c.JSON(http.StatusUnprocessableEntity, nil)
	}

	err = CriarTransacaoDB(*transacao)
	if err != nil {
		return err
	}

	responseJSON := map[string]interface{}{
		"limite": cliente.Limite,
		"saldo":  novoSaldo,
	}
	return c.JSON(http.StatusOK, responseJSON)
}
