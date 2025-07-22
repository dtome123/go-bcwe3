package contract

import (
	"context"

	"github.com/dtome123/go-bcwe3/eth/types"
)

type Contract interface {
	Call(ctx context.Context, method string, params ...any) ([]any, error)
	Transact(ctx context.Context, method string, privateKey string, params ...any) (*types.Tx, error)
	CallViewFunction(method string, params ...interface{}) (ViewResults, error)
}
