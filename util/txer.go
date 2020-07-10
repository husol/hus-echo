package util

import (
	"context"

	"github.com/jinzhu/gorm"
)

var keyEnable = "mysql_tx_enable"
var keyTx = "mysql_tx"

// nolint
func TxBegin(ctx context.Context, getClient func(ctx context.Context) *gorm.DB) context.Context {
	db := getClient(ctx)
	tx := db.Begin()
	ctx = SetTx(ctx, tx)

	ctx = context.WithValue(ctx, keyEnable, true)

	return ctx
}

// nolint
func TxEnd(ctx context.Context, txFunc func() error) (context.Context, error) {
	var err error

	tx := GetTx(ctx)

	defer func() {
		p := recover()

		switch {
		case p != nil:
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		case err != nil:
			tx.Rollback() // err is non-nil; don't change it
		default:
			err = tx.Commit().Error // if Commit returns error update err with commit err
		}
	}()

	err = txFunc()
	ctx = context.WithValue(ctx, keyEnable, false)

	return ctx, err
}

func IsEnableTx(ctx context.Context) bool {
	txEnable, ok := ctx.Value(keyEnable).(bool)

	return ok && txEnable
}

func GetTx(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(keyTx).(*gorm.DB)
	if !ok {
		return nil
	}

	return tx
}

// nolint
func SetTx(ctx context.Context, tx *gorm.DB) context.Context {
	ctx = context.WithValue(ctx, keyTx, tx)

	return ctx
}
