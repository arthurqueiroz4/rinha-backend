package impl

import (
	"context"
	"github.com/arthurqueiroz4/rinha-de-backend/types"
	"github.com/jackc/pgx/v5/pgxpool"
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
	row := db.QueryRow(context.Background(), "SELECT * FROM create_transaction_and_update_consumer($1, $2, $3, $4)", t.ConsumerID, t.Type, t.Value, t.Description)
	tr := &types.TransactionResponse{}
	err := row.Scan(&tr.Balance, &tr.Limit)
	return tr, err
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
