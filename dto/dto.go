package dto

type TransacaoDTO struct {
	Valor     int    `json:"valor"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
}
