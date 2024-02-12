package repo

import (
	"context"
	"github.com/arthurqueiroz4/rinha-de-backend/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransacaoRepository struct {
	db                *pgxpool.Pool
	clienteRepository *ClienteRepository
}

func NewTransacaoRepository(db *pgxpool.Pool, clienteRepository *ClienteRepository) *TransacaoRepository {
	return &TransacaoRepository{
		db:                db,
		clienteRepository: clienteRepository,
	}
}

func (r *TransacaoRepository) GetLastTransacoes(idCliente string) ([]types.Transacao, error) {
	cliente, err := r.clienteRepository.GetCliente(idCliente)
	if err != nil {
		return nil, err
	}
	transacoes := make([]types.Transacao, 0)
	rows, _ := r.db.Query(context.Background(), "SELECT * FROM transacoes t WHERE cliente_id = $1 ORDER BY t.id DESC LIMIT 10", cliente.ID)
	for rows.Next() {
		transacao := types.Transacao{}
		err := rows.Scan(&transacao.ID, &transacao.ClienteID, &transacao.Tipo, &transacao.Valor, &transacao.Descricao, &transacao.RealizadoEm)
		if err != nil {
			return nil, err
		}
		transacoes = append(transacoes, transacao)
	}

	return transacoes, nil
}

func (r *TransacaoRepository) CriarTransacao(transacao types.Transacao) error {
	tx, _ := r.db.Begin(context.Background())
	batch := pgx.Batch{}
	batch.Queue("UPDATE clientes SET saldo = saldo + $1 WHERE id = $2", transacao.Valor, transacao.ClienteID)
	batch.Queue("INSERT INTO transacoes (cliente_id, tipo, valor, descricao, realizado_em) VALUES ($1, $2, $3, $4, $5)",
		transacao.ClienteID, transacao.Tipo, transacao.Valor, transacao.Descricao, transacao.RealizadoEm)

	br := tx.SendBatch(context.Background(), &batch)
	_, err := br.Exec()
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	tx.Commit(context.Background())
	return nil
}
