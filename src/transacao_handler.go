package main

import (
	"github.com/labstack/echo"
	"net/http"
	"strconv"
	"sync"
)

func CriarTransacao(c echo.Context) error {
	id := c.Param("id")
	if idInt, _ := strconv.Atoi(id); idInt > 5 || idInt < 1 {
		return c.JSON(404, nil)
	}

	transacao := Transacao{}
	transacao.ClienteID, _ = strconv.Atoi(id)
	if err := c.Bind(&transacao); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, nil)
	}

	if transacao.Tipo != "d" && transacao.Tipo != "c" {
		return c.JSON(http.StatusUnprocessableEntity, nil)
	}

	if descLen := len(transacao.Descricao); descLen < 1 || descLen > 10 {
		return c.JSON(http.StatusUnprocessableEntity, nil)
	}

	m := sync.Mutex{}
	m.Lock()
	if err := CriarTransacaoDB(transacao); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, nil)
	}
	clienteDB, _ := PegarClienteDB(id)
	responseJSON := map[string]interface{}{
		"limite": clienteDB.Limite,
		"saldo":  clienteDB.Saldo,
	}
	m.Unlock()
	return c.JSON(http.StatusOK, responseJSON)
}
