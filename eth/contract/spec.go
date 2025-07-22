package contract

import (
	"context"

	"github.com/dtome123/go-bcwe3/eth/types"
)

type Contract interface {
	Call(ctx context.Context, method string, params ...any) (ContractResults, error)
	Transact(ctx context.Context, method string, privateKey string, params ...any) (*types.Tx, error)
	CallViewFunction(method string, params ...interface{}) (ContractResults, error)
	ListenContractEvent(eventName string, eventPrototype any, unpackFunc func(vLog types.Log, event interface{}) error, handleFunc func(event interface{})) error
}
