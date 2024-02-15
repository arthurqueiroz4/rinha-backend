package main

import (
	"context"
)

func PegarClienteDB(id string) (*Cliente, error) {
	cliente := &Cliente{}
	query := db.QueryRow(context.Background(), "SELECT * FROM clientes WHERE id = $1", id)
	err := query.Scan(&cliente.ID, &cliente.Limite, &cliente.Saldo)
	if err != nil {
		return nil, err
	}
	return cliente, nil
}
