package types

import "time"

type (
	Consumer struct {
		ID      int
		Balance int
		Limit   int
	}
	Transaction struct {
		ID          int       `json:"id"`
		ConsumerID  int       `json:"cliente_id"`
		Type        string    `json:"tipo"`
		Value       int       `json:"valor"`
		Description string    `json:"descricao"`
		CreatedAt   time.Time `json:"realizado_em"`
	}
)

type (
	Extract struct {
		Balance          Balance       `json:"saldo"`
		LastTransactions []Transaction `json:"ultimas_transacoes"`
	}
	Balance struct {
		Total       int       `json:"total"`
		ExtractDate time.Time `json:"data_extrato"`
		Limit       int       `json:"limite"`
	}
	TransactionResponse struct {
		ConsumerID int `json:"id"`
		Balance    int `json:"saldo"`
		Limit      int `json:"limite"`
	}
)
