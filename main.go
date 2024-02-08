package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arthurqueiroz4/rinha-go/config"
	"github.com/arthurqueiroz4/rinha-go/dto"
	"github.com/arthurqueiroz4/rinha-go/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"strconv"
)

func main() {
	log.Println("Starting the application...")
	env := &config.Env{}
	env.LoadEnv()
	env.ConnectDB()
	defer env.CloseDB()
	db := env.DB
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/clientes/{id}/transacoes", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if idInt, _ := strconv.Atoi(id); idInt > 5 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var transacaoDTO dto.TransacaoDTO

		err := json.NewDecoder(r.Body).Decode(&transacaoDTO)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if transacaoDTO.Tipo != "d" && transacaoDTO.Tipo != "c" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		descLength := len(transacaoDTO.Descricao)
		if descLength < 1 || descLength > 10 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var cliente model.Cliente
		err = db.QueryRow(context.Background(), "SELECT c.id, c.saldo, c.limite FROM clientes c WHERE id = $1", id).Scan(&cliente.ID, &cliente.Saldo, &cliente.Limite)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if transacaoDTO.Tipo == "d" {
			cliente.Saldo -= transacaoDTO.Valor
		} else {
			cliente.Saldo += transacaoDTO.Valor
		}

		abs := func(i int) int {
			if i < 0 {
				return -i
			}
			return i
		}

		if abs(cliente.Saldo) > cliente.Limite {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		db.QueryRow(context.Background(), "INSERT INTO transacoes (valor, tipo, descricao, cliente_id) VALUES ($1, $2, $3, $4)",
			transacaoDTO.Valor, transacaoDTO.Tipo, transacaoDTO.Descricao, cliente.ID)
		db.QueryRow(context.Background(), "UPDATE clientes SET saldo = $1 WHERE id = $2", cliente.Saldo, cliente.ID)

		log.Println("Transação realizada com sucesso...")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"limite": cliente.Limite,
			"saldo":  cliente.Saldo,
		})
	})

	r.Get("/clientes/{id}/extrato", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if idInt, _ := strconv.Atoi(id); idInt > 5 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		rows, _ := db.Query(context.Background(), "SELECT * FROM transacoes t WHERE cliente_id = $1", id)
		defer rows.Close()

		var transacoes []model.Transacao
		for rows.Next() {
			var transacao model.Transacao
			err := rows.Scan(&transacao.ID, &transacao.ClienteID, &transacao.Tipo, &transacao.Valor, &transacao.Realizado)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			transacoes = append(transacoes, transacao)
		}

		if err := rows.Err(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var cliente model.Cliente
		db.QueryRow(context.Background(), "SELECT saldo FROM clientes WHERE id = $1", id).Scan(&cliente)

		w.WriteHeader(http.StatusOK)
	})
	http.ListenAndServe(env.AppPort, r)
}
