package impl

import (
	"context"
	"errors"
	"github.com/arthurqueiroz4/rinha-de-backend/types"
	"github.com/jackc/pgx/v5/pgxpool"
	"math"
)

type Postgres struct{}

var db *pgxpool.Pool

func NewPostgres(dsn string) *Postgres {
	var err error
	db, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		panic("Connect database error")
	}
	return &Postgres{}
}

func (p Postgres) GetConsumer(id int) (*types.Consumer, error) {
	row := db.QueryRow(context.Background(), "SELECT id, balance, bound FROM consumers WHERE id = $1")
	consumer := &types.Consumer{}
	err := row.Scan(&consumer.ID, &consumer.Balance, &consumer.Limit)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

func (p Postgres) CreateTransationAndUpdateConsumer(t *types.Transaction) (*types.TransactionResponse, error) {
	tx, _ := db.Begin(context.Background())
	defer tx.Rollback(context.Background())
	row := tx.QueryRow(context.Background(), "SELECT id, balance, bound FROM consumers WHERE id = $1 FOR UPDATE", t.ConsumerID)
	c := &types.Consumer{}
	err := row.Scan(&c.ID, &c.Balance, &c.Limit)
	if err != nil {
		return nil, err
	}
	if t.Type == "d" && float64(c.Limit) < math.Abs(float64(c.Balance-t.Value)) {
		return nil, errors.New("limit exceeded")
	}

	tx.Exec(context.Background(),
		`INSERT INTO transactions (consumer_id, type, value, description, created_at) 
			 VALUES ($1, $2, $3, $4, $5)`, t.ConsumerID, t.Type, t.Value, t.Description, t.CreatedAt)
	if t.Type == "d" {
		tx.Exec(context.Background(), `UPDATE consumers SET balance = balance - $1 WHERE id = $2`, t.Value, t.ConsumerID)
		tx.Commit(context.Background())
		return &types.TransactionResponse{ConsumerID: t.ConsumerID, Balance: c.Balance - t.Value, Limit: c.Limit}, nil
	} else {
		tx.Exec(context.Background(), `UPDATE consumers SET balance = balance + $1 WHERE id = $2`, t.Value, t.ConsumerID)
		tx.Commit(context.Background())
		return &types.TransactionResponse{ConsumerID: t.ConsumerID, Balance: c.Balance + t.Value, Limit: c.Limit}, nil
	}
}

func (p *Postgres) GetExtract(id int) (*types.Extract, error) {
	q, err := db.Query(context.Background(), "SELECT * FROM transactions WHERE consumer_id = $1 ORDER BY id DESC LIMIT 10", id)
	if err != nil {
		return nil, err
	}
	extract := &types.Extract{}
	ts := make([]types.Transaction, 0)
	for q.Next() {
		t := types.Transaction{}
		err := q.Scan(&t.ID, &t.ConsumerID, &t.Type, &t.Value, &t.Description, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		ts = append(ts, t)
	}
	extract.LastTransactions = ts

	row := db.QueryRow(context.Background(), "SELECT balance, bound, now() FROM consumers WHERE id = $1", id)
	b := types.Balance{}
	err = row.Scan(&b.Total, &b.Limit, &b.ExtractDate)
	if err != nil {
		return nil, err
	}
	extract.Balance = b
	return extract, nil
}
