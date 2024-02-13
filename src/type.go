package main

import "time"

type (
	Cliente struct {
		ID     int `json:"id"`
		Limite int `json:"limite"`
		Saldo  int `json:"saldo"`
	}

	Transacao struct {
		ID          int       `json:"id"`
		ClienteID   int       `json:"cliente_id"`
		Tipo        string    `json:"tipo"`
		Valor       int       `json:"valor"`
		Descricao   string    `json:"descricao"`
		RealizadoEm time.Time `json:"realizado_em"`
	}

	Extrato struct {
		Saldo             Saldo       `json:"saldo"`
		UltimasTransacoes []Transacao `json:"ultimas_transacoes"`
	}

	Saldo struct {
		Total       int       `json:"total"`
		DataExtrato time.Time `json:"data_extrato"`
		Limite      int       `json:"limite"`
	}
)
