package main

import (
	"github.com/arthurqueiroz4/rinha-de-backend/api"
	"github.com/arthurqueiroz4/rinha-de-backend/storage/impl"
	"log"
	"os"
)

func main() {
	port := os.Getenv("APP_PORT")

	postgres := impl.NewPostgres(os.Getenv("DATABASE_URL"))

	s := api.NewServer(port, postgres)
	log.Print("APLICAÇÃO RODANDO NA PORTA", port)
	log.Fatal(s.Start())
}
