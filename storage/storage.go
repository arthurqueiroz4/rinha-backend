package storage

import "github.com/arthurqueiroz4/rinha-de-backend/types"

type Storage interface {
	GetConsumer(id string) (*types.Consumer, error)
	CreateTransation()
}
