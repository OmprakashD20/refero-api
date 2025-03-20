package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/OmprakashD20/refero-api/repository"
)

type Store struct {
	conn *pgxpool.Pool
	db   *repository.Queries
}

func NewTransactionStore(conn *pgxpool.Pool) *Store {
	return &Store{conn: conn, db: repository.New(conn)}
}

func (s *Store) Exec(ctx context.Context, fn func(q *repository.Queries) error) error {
	txn, err := s.conn.Begin(ctx)
	if err != nil {
		return err
	}

	q := s.db.WithTx(txn)

	err = fn(q)
	if err != nil {
		if rollbackErr := txn.Rollback(ctx); rollbackErr != nil {
			return fmt.Errorf("txn err: %v, rollback err: %v", err, rollbackErr)
		}

		return err
	}

	return txn.Commit(ctx)

}
