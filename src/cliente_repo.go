package main

import (
	"context"
)

func PegarClienteDB(id string) (*Cliente, error) {
	cliente, found := cache[id]

	if found {
		return cliente, nil
	}

	cliente = &Cliente{}
	query := db.QueryRow(context.Background(), "SELECT * FROM clientes WHERE id = $1", id)
	err := query.Scan(&cliente.ID, &cliente.Limite, &cliente.Saldo)
	if err != nil {
		return nil, err
	}
	cache[id] = cliente
	return cliente, nil
}
