package api

import (
	"encoding/json"
	"errors"
	"github.com/arthurqueiroz4/rinha-de-backend/storage/interface"
	"github.com/arthurqueiroz4/rinha-de-backend/types"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Server struct {
	port  string
	store _interface.Storage
}

func NewServer(port string, store _interface.Storage) *Server {
	return &Server{port: port, store: store}
}

func (s *Server) Start() error {
	r := mux.NewRouter()
	r.HandleFunc("/clientes/{id}/extrato", s.handleGetExtract)
	r.HandleFunc("/clientes/{id}/transacoes", s.handleCreateTransaction)
	return http.ListenAndServe(s.port, r)
}

func (s *Server) handleGetExtract(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	extract, err := s.store.GetExtract(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(extract)
}

func (s *Server) handleCreateTransaction(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	t := &types.Transaction{}
	err = json.NewDecoder(r.Body).Decode(t)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if l := len(t.Description); l > 10 || l < 1 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if t.Type != "c" && t.Type != "d" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	t.ConsumerID = id
	response, err := s.store.CreateTransationAndUpdateConsumer(t)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func getId(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return 0, errors.New("Id não encontrado no path")
	}
	return verifyId(id)
}

func verifyId(id string) (int, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil || idInt > 5 || idInt < 1 {
		return 0, errors.New("Invalid ID")
	}
	return idInt, nil
}
