package listener

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Listener listens to Ethereum contract events.
type Listener struct {
	url string
}

// NewListener creates a new listener.
func NewListener(url string) *Listener {
	return &Listener{url: url}
}

// EventCallback handles decoded events.
type EventCallback func(log types.Log)

// ListenEvents subscribes to events (all if eventNames empty).
// It uses a buffered channel + auto reconnect + keepalive ping.
func (l *Listener) ListenEvents(
	ctx context.Context,
	contractAddr string,
	parsedABI abi.ABI,
	eventNames []string,
	cb EventCallback,
) {

	go func(ctx context.Context) {
		keepaliveTicker := time.NewTicker(10 * time.Second)
		defer keepaliveTicker.Stop()

		client, err := ethclient.DialContext(ctx, l.url)
		if err != nil {
			fmt.Printf("[listener] failed to dial: %v\n", err)
			return
		}

		for {
			select {
			case <-ctx.Done():
				fmt.Println("[listener] context canceled, stopping listener")
				return
			default:
			}

			// Build topics
			var topics [][]common.Hash
			if len(eventNames) > 0 {
				var sigs []common.Hash
				for _, name := range eventNames {
					if ev, ok := parsedABI.Events[name]; ok {
						sigs = append(sigs, ev.ID)
					} else {
						fmt.Printf("[listener] event %s not found in ABI\n", name)
					}
				}
				topics = [][]common.Hash{sigs}
			}

			latestBlock, err := client.BlockNumber(ctx)
			if err != nil {
				fmt.Printf("[listener] failed to get latest block: %v\n", err)
				return
			}

			query := ethereum.FilterQuery{
				Addresses: []common.Address{common.HexToAddress(contractAddr)},
				Topics:    topics,
				FromBlock: big.NewInt(int64(latestBlock)),
			}

			// Buffered channel to avoid blocking
			logsChan := make(chan types.Log, 1000)

			sub, err := client.SubscribeFilterLogs(ctx, query, logsChan)
			if err != nil {
				fmt.Printf("[listener] subscribe failed: %v, retrying...\n", err)
				time.Sleep(time.Second) // simple retry delay
				continue
			}

			fmt.Printf("[listener] subscription established for %s\n", contractAddr)

		subLoop:
			for {
				select {
				case <-ctx.Done():
					sub.Unsubscribe()
					fmt.Println("[listener] context canceled, stopping listener")
					return
				case err := <-sub.Err():
					fmt.Printf("[listener] subscription error: %v, reconnecting...\n", err)
					sub.Unsubscribe()

					// check specific websocket close code

					if shouldRedial(err) {
						fmt.Printf("[listener] websocket closed (%v), re-dialing...\n", err)

						// new dial
						newClient, dialErr := retryDial(ctx, l.url)
						if dialErr != nil {
							fmt.Printf("[listener] failed to retry dial: %v\n", dialErr)
							return
						}

						// reconnect
						client = newClient
					}

					break subLoop
				case vLog := <-logsChan:
					fmt.Printf("[listener] contract received event, topic length: %d\n", len(vLog.Topics))
					if len(vLog.Topics) > 0 {
						fmt.Println("[listener] contract callback")
						cb(vLog)
					}
				case <-keepaliveTicker.C:
					// simple keepalive ping
					_, err := client.BlockNumber(ctx)
					if err != nil {
						fmt.Printf("[listener] keepalive failed: %v\n", err)
					}
				}
			}

			// small delay before reconnect
			time.Sleep(time.Second)
		}
	}(ctx)
}

func shouldRedial(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	switch {
	case strings.Contains(msg, "1000"), // Normal closure
		strings.Contains(msg, "1001"), // Going away
		strings.Contains(msg, "1011"), // Internal error
		strings.Contains(msg, "1012"), // Service restart
		strings.Contains(msg, "1013"): // Try again later
		return true
	default:
		return false
	}
}

func retryDial(ctx context.Context, url string) (*ethclient.Client, error) {
	for {
		client, err := ethclient.Dial(url)
		if err == nil {
			return client, nil
		}
		select {
		case <-ctx.Done():
			return nil, err
		case <-time.After(time.Second):
		}
	}
}
