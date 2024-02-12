package handler

import (
	"github.com/arthurqueiroz4/rinha-de-backend/repo"
	"github.com/arthurqueiroz4/rinha-de-backend/types"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

func CriarTransacao(c echo.Context) error {
	tr := c.Get("tr").(*repo.TransacaoRepository)
	cr := c.Get("cr").(*repo.ClienteRepository)
	id := c.Param("id")
	if idInt, _ := strconv.Atoi(id); idInt > 5 || idInt < 1 {
		return c.JSON(404, nil)
	}

	transacao := &types.Transacao{}
	if err := c.Bind(transacao); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, nil)
	}

	if transacao.Tipo != "d" && transacao.Tipo != "c" {
		return c.JSON(http.StatusUnprocessableEntity, nil)
	} else if transacao.Tipo == "d" {
		transacao.Valor = -transacao.Valor
	}

	if descLen := len(transacao.Descricao); descLen < 1 || descLen > 10 {
		return c.JSON(http.StatusUnprocessableEntity, nil)
	}

	cliente, err := cr.GetCliente(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "deu erro no GetCliente")
	}
	if (cliente.Saldo + transacao.Valor) < -cliente.Limite {
		return c.JSON(http.StatusUnprocessableEntity, nil)
	}

	err = tr.CriarTransacao(*transacao)
	if err != nil {
		return err
	}
	cliente.Saldo += transacao.Valor

	responseJSON := map[string]interface{}{
		"limite": cliente.Limite,
		"saldo":  cliente.Saldo,
	}
	return c.JSON(http.StatusOK, responseJSON)
}
