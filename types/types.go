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
	CreateTransaction struct {
		ConsumerID int
		Balance    int
		Type       string
	}
)

func (ct *CreateTransaction) NewCreateTransaction(consumerID, balance int, typeTransaction string) *CreateTransaction {
	if typeTransaction == "d" {
		return &CreateTransaction{Balance: balance * -1, Type: typeTransaction, ConsumerID: consumerID}
	}
	return &CreateTransaction{Balance: balance, Type: typeTransaction, ConsumerID: consumerID}
}
