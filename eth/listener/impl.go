package listener

import (
	"context"
	"fmt"

	"github.com/dtome123/go-bcwe3/eth/provider"
	"github.com/dtome123/go-bcwe3/eth/types"

	goethTypes "github.com/ethereum/go-ethereum/core/types"
)

type impl struct {
	client provider.Provider
}

func NewListener(
	client provider.Provider,
) Listener {
	return &impl{
		client: client,
	}
}

func (l *impl) ListenEventBlock(handle func(block *types.Block), errorChan chan<- error) {
	headers := make(chan *goethTypes.Header)
	sub, err := l.client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		errorChan <- fmt.Errorf("error subscribing to new blocks: %w", err)
		return
	}

	for {
		select {
		case err := <-sub.Err():
			if err != nil {
				errorChan <- fmt.Errorf("subscription error: %w", err)
			}
		case header := <-headers:

			block, err := l.client.BlockByNumber(context.Background(), header.Number)
			if err != nil {
				errorChan <- fmt.Errorf("error fetching block: %w", err)
				continue
			}

			handle(block)

		}
	}
}
