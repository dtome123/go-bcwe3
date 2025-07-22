package contract

import (
	"context"

	"github.com/dtome123/go-bcwe3/eth/types"
)

type Contract interface {
	Transact(ctx context.Context, method string, privateKey string, params ...any) (*types.Tx, error)
	Call(ctx context.Context, method string, params ...interface{}) (ContractResults, error)
	ListenEvent(ctx context.Context, eventName string, eventPrototype any, unpackFunc func(vLog types.Log, event interface{}) error, handleFunc func(event interface{})) error
}
