package model

import "time"

type Transacao struct {
	ID        int       `json:"id"`
	Valor     int       `json:"valor"`
	Tipo      string    `json:"tipo"`
	Descricao string    `json:"descricao"`
	ClienteID int       `json:"cliente_id"`
	Realizado time.Time `json:"realizada_em"`
}

type Cliente struct {
	ID     int `json:"id"`
	Saldo  int `json:"saldo"`
	Limite int `json:"limite"`
}
