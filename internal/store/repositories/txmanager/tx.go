package txmanager

import (
	"context"
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
