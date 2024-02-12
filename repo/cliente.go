package repo

import (
	"context"
	"github.com/arthurqueiroz4/rinha-de-backend/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClienteRepository struct {
	db *pgxpool.Pool
}

var cache map[string]*types.Cliente

func NewClienteRepository(db *pgxpool.Pool) *ClienteRepository {
	cache = make(map[string]*types.Cliente)
	return &ClienteRepository{
		db: db,
	}
}

func (r *ClienteRepository) GetCliente(id string) (*types.Cliente, error) {
	cliente, found := cache[id]

	if found {
		return cliente, nil
	}

	cliente = &types.Cliente{}
	query := r.db.QueryRow(context.Background(), "SELECT * FROM clientes WHERE id = $1", id)
	err := query.Scan(&cliente.ID, &cliente.Limite, &cliente.Saldo)
	if err != nil {
		return nil, err
	}
	cache[id] = cliente
	return cliente, nil
}
