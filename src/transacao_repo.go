package main

import (
	"context"
)

func PegarUltimasTransacoesDB(idCliente string) ([]Transacao, error) {
	cliente, err := PegarClienteDB(idCliente)
	if err != nil {
		return nil, err
	}
	transacoes := make([]Transacao, 0)
	rows, _ := db.Query(context.Background(), "SELECT * FROM transacoes t WHERE cliente_id = $1 ORDER BY t.id DESC LIMIT 10", cliente.ID)
	for rows.Next() {
		transacao := Transacao{}
		err := rows.Scan(&transacao.ID, &transacao.ClienteID, &transacao.Tipo, &transacao.Valor, &transacao.Descricao, &transacao.RealizadoEm)
		if err != nil {
			return nil, err
		}
		transacoes = append(transacoes, transacao)
	}

	return transacoes, nil
}

func CriarTransacaoDB(transacao Transacao) error {
	var err error
	_, err = db.Exec(context.Background(), "UPDATE clientes SET saldo = saldo + $1 WHERE id = $2", transacao.Valor, transacao.ClienteID)
	if err != nil {
		return err
	}
	_, err = db.Exec(context.Background(), "INSERT INTO transacoes (cliente_id, tipo, valor, descricao, realizado_em) VALUES ($1, $2, $3, $4, $5)",
		transacao.ClienteID, transacao.Tipo, transacao.Valor, transacao.Descricao, transacao.RealizadoEm)
	if err != nil {
		return err
	}

	return err
}
