package main

import (
	"context"
	"errors"
)

func PegarUltimasTransacoesDB(idCliente string) ([]Transacao, error) {
	transacoes := make([]Transacao, 0)
	rows, _ := db.Query(context.Background(), "SELECT * FROM transacoes t WHERE cliente_id = $1 ORDER BY t.id DESC LIMIT 10", idCliente)
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
	tx, err := db.Begin(context.Background())
	defer tx.Rollback(context.Background())
	if err != nil {
		return err
	}
	var saldo, limite int
	err = tx.QueryRow(context.Background(), "SELECT saldo, limite FROM clientes WHERE id = $1", transacao.ClienteID).Scan(&saldo, &limite)
	if err != nil {
		return err
	}
	if transacao.Tipo == "d" {
		if saldo-transacao.Valor < -limite {
			return errors.New("limite excedido")
		}
		exec, err := tx.Exec(context.Background(), "UPDATE clientes SET saldo = saldo - $1 WHERE id = $2 AND saldo - $1 > -limite", transacao.Valor, transacao.ClienteID)
		if err != nil {
			return err
		}
		if exec.RowsAffected() == 0 {
			return errors.New("não atualizou saldo")
		}
	} else {
		exec, err := tx.Exec(context.Background(), "UPDATE clientes SET saldo = saldo + $1 WHERE id = $2", transacao.Valor, transacao.ClienteID)
		if err != nil {
			return err
		}
		if exec.RowsAffected() == 0 {
			return errors.New("não atualizou saldo")
		}
	}

	exec, err := tx.Exec(context.Background(), "INSERT INTO transacoes (cliente_id, tipo, valor, descricao, realizado_em) VALUES ($1, $2, $3, $4, $5)",
		transacao.ClienteID, transacao.Tipo, transacao.Valor, transacao.Descricao, transacao.RealizadoEm)
	if err != nil {
		return err
	}
	if exec.RowsAffected() == 0 {
		return errors.New("não inseriu transacao")
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}
