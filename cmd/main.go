package main

import (
	"context"
	"github.com/arthurqueiroz4/rinha-de-backend/api/server"
	"github.com/arthurqueiroz4/rinha-de-backend/repo"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

func main() {
	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	cr := repo.NewClienteRepository(db)
	tr := repo.NewTransacaoRepository(db, cr)

	server.NewServer(cr, tr)
}
