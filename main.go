package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/arthurqueiroz4/rinha-go/config"
	"github.com/arthurqueiroz4/rinha-go/dto"
	"github.com/arthurqueiroz4/rinha-go/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

		if idInt, _ := strconv.Atoi(id); idInt > 5 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var cliente model.Cliente
		db.Raw("SELECT * FROM clientes c WHERE c.id = ?", id).Scan(&cliente)

		if cliente.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		abs := func(i int) int {
			if i < 0 {
				return -i
			}
			return i
		}
		if ((abs(cliente.Saldo)) + transacaoDTO.Valor) > cliente.Limite {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		if err := db.Raw("INSERT INTO transacoes (valor, tipo, descricao, cliente_id) VALUES (?, ?, ?, ?)",
			transacaoDTO.Valor, transacaoDTO.Tipo, transacaoDTO.Descricao, cliente.ID).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println("Transação registrada com sucesso...")
		if err := db.Raw("UPDATE clientes SET saldo = saldo - ? WHERE id = ?", transacaoDTO.Valor, cliente.ID).Error; err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println("Transação realizada com sucesso...")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"limite": cliente.Limite,
			"saldo":  cliente.Saldo - transacaoDTO.Valor,
		})
	})
	http.ListenAndServe(env.AppPort, r)
}
