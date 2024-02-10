package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	dbPool        *pgxpool.Pool
	cache         map[string]Extrato
	cacheMux      sync.Mutex
	cacheClientes map[string]Cliente
)

func main() {
	var err error
	databaseUrl := "postgres://postgres:postgres@db:5432/rinha"
	dbPool, err = pgxpool.New(context.Background(), databaseUrl)
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
	ClienteID int       `json:"cliente_id"`
	Tipo      string    `json:"tipo"`
	Valor     int       `json:"valor"`
	Descricao string    `json:"descricao"`
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
	transacao := new(Transacao)
	validation := make(chan error)

	go func(validation chan error) {
		if idInt, _ := strconv.Atoi(id); idInt > 5 || idInt < 1 {
			validation <- c.JSON(http.StatusNotFound, nil)
			return
		}

		if err := c.Bind(transacao); err != nil {
			validation <- c.JSON(http.StatusUnprocessableEntity, nil)
			return
		}

		if transacao.Tipo != "d" && transacao.Tipo != "c" {
			validation <- c.JSON(http.StatusUnprocessableEntity, nil)
			return
		}

		descLen := len(transacao.Descricao)
		if descLen < 1 || descLen > 10 {
			validation <- c.JSON(http.StatusUnprocessableEntity, nil)
			return
		}
		validation <- nil
	}(validation)

	cliente := new(Cliente)
	err := dbPool.QueryRow(context.Background(), "SELECT * FROM clientes c WHERE id = $1", id).Scan(&cliente.ID, &cliente.Limite, &cliente.Saldo)

	var novoSaldo int
	if transacao.Tipo == "d" {
		if (cliente.Saldo - transacao.Valor) < -cliente.Limite {
			return c.JSON(http.StatusUnprocessableEntity, nil)
		}
		novoSaldo = cliente.Saldo - transacao.Valor
	} else {
		novoSaldo = cliente.Saldo + transacao.Valor
	}

	tx, err := dbPool.Begin(context.Background())
	defer tx.Rollback(context.Background())

	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	batch := &pgx.Batch{}
	batch.Queue("INSERT INTO transacoes (valor, tipo, descricao, cliente_id, realizado) VALUES ($1, $2, $3, $4, $5)",
		transacao.Valor, transacao.Tipo, transacao.Descricao, id, time.Now())
	batch.Queue("UPDATE clientes SET saldo = $1 WHERE id = $2", novoSaldo, id)

	err = <-validation
	if err != nil {
		return err
	}

	br := tx.SendBatch(context.Background(), batch)
	_, err = br.Exec()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	err = br.Close()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	responseJSON := map[string]interface{}{
		"limite": cliente.Limite,
		"saldo":  novoSaldo,
	}
	return c.JSON(http.StatusOK, responseJSON)
}

func pegarExtrato(c echo.Context) error {
	id := c.Param("id")
	if idInt, _ := strconv.Atoi(id); idInt > 5 || idInt < 1 {
		return c.JSON(http.StatusNotFound, nil)
	}

	//cacheMux.Lock()
	//cachedExtrato, found := cache[id]
	//cacheMux.Unlock()
	//if found {
	//	return c.JSON(http.StatusOK, cachedExtrato)
	//}

	rows, _ := dbPool.Query(context.Background(), "SELECT * FROM transacoes t WHERE cliente_id = $1 ORDER BY t.id DESC LIMIT 10", id)
	ultimasTransacoes, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Transacao])
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	cliente := new(Cliente)
	var dataExtrato time.Time
	dbPool.QueryRow(context.Background(), "SELECT c.id, c.saldo, c.limite, now() FROM clientes c WHERE id = $1", id).Scan(&cliente.ID, &cliente.Saldo, &cliente.Limite, &dataExtrato)

	extrato := Extrato{
		Saldo: Saldo{
			Total:       cliente.Saldo,
			DataExtrato: dataExtrato,
			Limite:      cliente.Limite,
		},
		UltimasTransacoes: ultimasTransacoes,
	}

	//cacheMux.Lock()
	//cache[id] = extrato
	//cacheMux.Unlock()

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
