package listener

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/dtome123/go-bcwe3/eth/provider"
	"github.com/dtome123/go-bcwe3/eth/types"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
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

func (l *impl) ListenBlock(handleFunc func(block *types.Block), errorChan chan<- error) {
	headers := make(chan *types.Header)
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

			handleFunc(block)

		}
	}
}

func (l *impl) ListenContractEvent(
	contractAddress string,
	eventName string,
	eventPrototype any,
	unpackFunc func(vLog types.Log, event interface{}) error,
	handleFunc func(event interface{}),
) error {

	address := common.HexToAddress(contractAddress)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{address},
	}

	logs := make(chan types.Log)
	sub, err := l.client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return fmt.Errorf("failed to subscribe to logs: %w", err)
	}

	go func() {
		for {
			select {
			case err := <-sub.Err():
				log.Println("Subscription error:", err)
			case vLog := <-logs:
				// Clone prototype
				eventCopy := reflect.New(reflect.TypeOf(eventPrototype).Elem()).Interface()

				err := unpackFunc(vLog, eventCopy)
				if err != nil {
					log.Println("Unpack failed:", err)
					continue
				}

				handleFunc(eventCopy)
			}
		}
	}()

	return nil
}
