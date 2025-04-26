package listener

import (
	"github.com/dtome123/go-bcwe3/eth/types"
)

type Listener interface {
	ListenBlock(handleFunc func(block *types.Block), errorChan chan<- error)
	ListenContractEvent(contractAddress string, eventName string, eventPrototype any, unpackFunc func(vLog types.Log, event interface{}) error, handleFunc func(event any)) error
}
