package listener

import "github.com/dtome123/go-bcwe3/eth/types"

type Listener interface {
	ListenEventBlock(handle func(block *types.Block), errorChan chan<- error)
}
