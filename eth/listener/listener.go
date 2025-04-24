package listener

import (
	"context"
	"fmt"
	"time"

	"github.com/dtome123/go-bcwe3/eth/client"
	"github.com/dtome123/go-bcwe3/eth/types"

	goethTypes "github.com/ethereum/go-ethereum/core/types"
)

type Listener interface {
	Start(handle func(block *types.Block), errorChan chan<- error, onTick func() error)
}

type impl struct {
	client client.Client
}

func NewListener(client client.Client) Listener {
	return &impl{
		client: client,
	}
}

func (l *impl) Start(handle func(block *types.Block), errorChan chan<- error, onTick func() error) {
	headers := make(chan *goethTypes.Header)
	sub, err := l.client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		errorChan <- fmt.Errorf("error subscribing to new blocks: %w", err)
		return
	}

	// Set up a timeout and ticker
	timeoutDuration := time.Minute // Set timeout duration, e.g., 1 minute
	timeout := time.After(timeoutDuration)
	ticker := time.NewTicker(time.Second * 10) // Check every 10 seconds
	defer ticker.Stop()

	for {
		select {
		case err := <-sub.Err():
			if err != nil {
				errorChan <- fmt.Errorf("subscription error: %w", err)
			}
		case header := <-headers:
			block, err := l.client.BlockByHash(context.Background(), header.Hash().Hex())
			if err != nil {
				errorChan <- fmt.Errorf("error fetching block: %w", err)
				continue
			}
			handle(block)
		case <-ticker.C:
			// Execute the custom onTick callback
			if onTick != nil {
				if err := onTick(); err != nil {
					errorChan <- fmt.Errorf("onTick error: %w", err)
				}
			}
		case <-timeout:
			errorChan <- fmt.Errorf("listener timed out after %v without new blocks", timeoutDuration)
			return // Stop the listener
		}
	}
}
