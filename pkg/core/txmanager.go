package core

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

type TxCtxKey string

const (
	TxCtxKeyCode TxCtxKey = "txCtxKey"
)

type TxManager interface {
	Do(context.Context, func(context.Context) error) error

	DoWithSettings(
		context.Context,
		TxSettings,
		func(context.Context) error,
	) error
}

type TxSettings interface {
	GetTxCtxKey() TxCtxKey
	GetIsolationLevel() string
}

type gormTxManager struct {
	db *gorm.DB
}

func NewGormTxManager(db *gorm.DB) *gormTxManager {
	return &gormTxManager{
		db: db,
	}
}

func (txm *gormTxManager) Do(ctx context.Context, f func(context.Context) error) error {
	return txm.DoWithSettings(ctx, DefaultGormTxSettings(), f)
}

func (txm *gormTxManager) DoWithSettings(ctx context.Context, txSettings TxSettings, f func(context.Context) error) error {
	if existing := txm.getTxFromCtx(ctx); existing != nil {
		return f(ctx)
	}

	tx := txm.db.Begin(&sql.TxOptions{
		Isolation: txm.mapIsolationLevel(txSettings.GetIsolationLevel()),
	})
	if err := tx.Error; err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	txCtx := context.WithValue(ctx, TxCtxKeyCode, tx)
	err := f(txCtx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (txm *gormTxManager) getTxFromCtx(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(TxCtxKeyCode).(*gorm.DB)
	if !ok {
		return nil
	}
	return tx
}

func (txm *gormTxManager) mapIsolationLevel(isolationLevel string) sql.IsolationLevel {
	switch isolationLevel {
	case "Default":
		return sql.LevelDefault
	case "ReadCommitted":
		return sql.LevelReadCommitted
	case "Serializable":
		return sql.LevelSerializable
	default:
		return sql.LevelReadCommitted
	}
}

type gormTxSettings struct {
	isolationLevel string
}

func NewGormTxSettings(isolationLevel string) *gormTxSettings {
	return &gormTxSettings{
		isolationLevel: isolationLevel,
	}
}

func DefaultGormTxSettings() *gormTxSettings {
	return &gormTxSettings{
		isolationLevel: "ReadCommitted",
	}
}

func (txs *gormTxSettings) GetTxCtxKey() TxCtxKey {
	return TxCtxKeyCode
}

func (txs *gormTxSettings) GetIsolationLevel() string {
	return txs.isolationLevel
}
