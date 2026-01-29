package db

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type inTransact struct {}
type txKeyStruct = struct{}

var TxKey = txKeyStruct{}
var InTransactKey = inTransact{}
const (
	InTransactValue = "1"
)

func RunInTx(ctx context.Context, db *sqlx.DB, txFn func(ctx context.Context, tx *sqlx.Tx) error) error {
	var err error
	tx, ok := ctx.Value(TxKey).(*sqlx.Tx)
	if ok {
		logrus.Debug("sqlx.Tx is found in  RunInTx")
		return txFn(context.WithValue(ctx, TxKey, tx), tx)
	}
	logrus.Debug("beginning new transaction")
	// トランザクション開始
	tx, err = db.Beginx()
	if err != nil {
		return err
	}
	// トランザクション実行
	err = txFn(context.WithValue(ctx, TxKey, tx), tx)
	if err != nil {
		// ロールバック
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	inTx, ok := ctx.Value(InTransactKey).(string)
	if ok && inTx == InTransactValue {
		return nil
	}

	// コミット
	return tx.Commit()
}
