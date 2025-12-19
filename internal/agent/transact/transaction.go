package transact

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

type txKeyType struct{}

var txKey txKeyType = txKeyType{}

type TransactionProvider struct {
	db *sqlx.DB
}

func NewTransactionProvider(db *sqlx.DB) *TransactionProvider {
	if db == nil {
		panic("db is nil")
	}
	return &TransactionProvider{
		db: db,
	}
}

func (t *TransactionProvider) Transact(ctx context.Context, txFn func(context.Context) error) error {
	tx, err := t.db.Beginx()
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, txKey, tx)
	err = txFn(ctx)
	if err != nil {
		// ロールバック
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}

func Transact(ctx context.Context, txFn func(*sqlx.Tx) error) error {
	txValue := ctx.Value(txKey)
	tx, ok := txValue.(*sqlx.Tx)
	if !ok {
		return errors.New("invalid context to transact")
	}
	return txFn(tx)
}
