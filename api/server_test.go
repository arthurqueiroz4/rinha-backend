package api

import (
	"github.com/arthurqueiroz4/rinha-de-backend/types"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type storageTest struct{}

func (s *storageTest) GetConsumer(id int) (*types.Consumer, error) {
	return &types.Consumer{}, nil
}

func (s *storageTest) CreateTransationAndUpdateConsumer(t *types.Transaction) error {
	return nil
}

func (s *storageTest) GetExtract(id int) (*types.Extract, error) {
	return &types.Extract{}, nil
}

func TestNotFound(t *testing.T) {
	s := Server{port: ":3000", store: &storageTest{}}
	go s.Start()

	resp, _ := http.Get("http://localhost:3000/clientes/6/extrato")
	assert.Equal(t, 404, resp.StatusCode)

	resp, _ = http.Get("http://localhost:3000/clientes/0/extrato")
	assert.Equal(t, 404, resp.StatusCode)

	resp, _ = http.Get("http://localhost:3000/clientes/2/extrato")
	assert.NotEqual(t, 404, resp.StatusCode)
}
