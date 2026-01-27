package transact

import (
	"context"

	"github.com/TokujouKaisenDonburi/optical-backend/pkg/db"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)
type TransactionProvider interface {
	Transact(ctx context.Context, txFn func(context.Context) error) error
}

type PsqlTransactionProvider struct {
	db *sqlx.DB
}

func NewPsqlTransactionProvider(db *sqlx.DB) *PsqlTransactionProvider {
	return &PsqlTransactionProvider{
		db: db,
	}
}

func (p *PsqlTransactionProvider) Transact(ctx context.Context, txFn func(context.Context) error) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, db.TxKey, tx)
	err = txFn(ctx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			logrus.WithError(err).Error("failed to rollback in trsansaction")
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		logrus.WithError(err).Error("failed to commit transaction")
	}
	logrus.Debug("success to commit transaction")
	return nil
}
