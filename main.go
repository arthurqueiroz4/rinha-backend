package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
)

var (
	dbPool   *pgxpool.Pool
	cache    map[string]Extrato
	cacheMux sync.Mutex
)

func main() {
	var err error
	databaseUrl := "postgres://postgres:postgres@db:5432/rinha"
	dbPool, err = pgxpool.Connect(context.Background(), databaseUrl)
	if err != nil {
		fmt.Println("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbPool.Close()
	if err := dbPool.Ping(context.Background()); err != nil {
		fmt.Println("Unable to ping database: %v\n", err)
		os.Exit(1)
	}

	cache = make(map[string]Extrato, 2000)

	e := echo.New()
	e.Use(logRequest)
	e.POST("/clientes/:id/transacoes", criarTransacao)
	e.GET("/clientes/:id/extrato", pegarExtrato)
	e.Logger.Fatal(e.Start(":8080"))
}

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

type Extrato struct {
	Saldo             Saldo       `json:"saldo"`
	UltimasTransacoes []Transacao `json:"ultimas_transacoes"`
}

type Saldo struct {
	Total       int       `json:"total"`
	DataExtrato time.Time `json:"data_extrato"`
	Limite      int       `json:"limite"`
}

func criarTransacao(c echo.Context) error {
	id := c.Param("id")
	if idInt, _ := strconv.Atoi(id); idInt > 5 || idInt < 1 {
		return c.JSON(http.StatusNotFound, nil)
	}

	transacao := new(Transacao)
	if err := c.Bind(transacao); err != nil {
		return err
	}

	if transacao.Tipo != "d" && transacao.Tipo != "c" {
		return c.JSON(http.StatusUnprocessableEntity, nil)
	}

	descLen := len(transacao.Descricao)
	if descLen < 1 || descLen > 10 {
		return c.JSON(http.StatusUnprocessableEntity, nil)
	}

	cliente := new(Cliente)
	err := dbPool.QueryRow(context.Background(), "SELECT * FROM clientes c WHERE id = $1", id).Scan(&cliente.ID, &cliente.Limite, &cliente.Saldo)

	if transacao.Tipo == "d" {
		if (cliente.Saldo - transacao.Valor) < -cliente.Limite {
			return c.JSON(http.StatusUnprocessableEntity, nil)
		}
	}

	_, err = dbPool.Exec(context.Background(), "INSERT INTO transacoes (valor, tipo, descricao, cliente_id, realizado) VALUES ($1, $2, $3, $4, $5)",
		transacao.Valor, transacao.Tipo, transacao.Descricao, id, time.Now())

	var novoSaldo int
	if transacao.Tipo == "d" {
		novoSaldo = cliente.Saldo - transacao.Valor
	} else {
		novoSaldo = cliente.Saldo + transacao.Valor
	}

	_, err = dbPool.Exec(context.Background(), "UPDATE clientes SET saldo = $1 WHERE id = $2", novoSaldo, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	responseJSON := map[string]interface{}{
		"limite": cliente.Limite,
		"saldo":  novoSaldo,
	}

	cacheMux.Lock()
	delete(cache, id)
	cacheMux.Unlock()

	return c.JSON(http.StatusOK, responseJSON)
}

func pegarExtrato(c echo.Context) error {
	id := c.Param("id")
	if idInt, _ := strconv.Atoi(id); idInt > 5 || idInt < 1 {
		return c.JSON(http.StatusNotFound, nil)
	}

	cacheMux.Lock()
	cachedExtrato, found := cache[id]
	cacheMux.Unlock()
	if found {
		return c.JSON(http.StatusOK, cachedExtrato)
	}

	querySelectTransacoes, _ := dbPool.Query(context.Background(), "SELECT * FROM transacoes t WHERE cliente_id = $1 ORDER BY t.realizado DESC LIMIT 10", id)

	defer querySelectTransacoes.Close()

	transacoes := make([]Transacao, 0)
	for querySelectTransacoes.Next() {
		transacao := Transacao{}
		querySelectTransacoes.Scan(&transacao.ID, &transacao.ClienteID, &transacao.Tipo, &transacao.Valor, &transacao.Descricao, &transacao.Realizado)
		transacoes = append(transacoes, transacao)
	}

	cliente := new(Cliente)
	dbPool.QueryRow(context.Background(), "SELECT c.id, c.saldo, c.limite FROM clientes c WHERE id = $1", id).Scan(&cliente.ID, &cliente.Saldo, &cliente.Limite)

	var dataExtrato time.Time
	if len(transacoes) > 0 {
		dataExtrato = transacoes[0].Realizado
	}

	extrato := Extrato{
		Saldo: Saldo{
			Total:       cliente.Saldo,
			DataExtrato: dataExtrato,
			Limite:      cliente.Limite,
		},
		UltimasTransacoes: transacoes,
	}

	cacheMux.Lock()
	cache[id] = extrato
	cacheMux.Unlock()

	return c.JSON(http.StatusOK, extrato)
}

func logRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		// Executar o manipulador de solicitação
		err := next(c)

		// Registrar detalhes da solicitação
		fmt.Printf(
			"%s %s - %s %s - %s - %d\n",
			c.Request().Method,
			c.Request().URL.Path,
			c.RealIP(),
			time.Since(start),
			c.Response().Status,
			c.Response().Size,
		)

		return err
	}
}
