package _interface

import "github.com/arthurqueiroz4/rinha-de-backend/types"

type Storage interface {
	GetConsumer(int) (*types.Consumer, error)
	CreateTransationAndUpdateConsumer(*types.Transaction) (*types.TransactionResponse, error)
	GetExtract(int) (*types.Extract, error)
}
