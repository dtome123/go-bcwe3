package listener

import (
	"github.com/dtome123/go-bcwe3/eth/types"
	goethTypes "github.com/ethereum/go-ethereum/core/types"
)

type Listener interface {
	ListenBlock(handleFunc func(block *types.Block), errorChan chan<- error)
	ListenContractEvent(contractAddress string, eventName string, eventPrototype any, unpackFunc func(vLog goethTypes.Log, event interface{}) error, handleFunc func(event any)) error
}
