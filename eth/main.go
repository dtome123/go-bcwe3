package eth

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Ethereum struct {
	Client *ethclient.Client
}

func NewEthereum(dsn string) *Ethereum {
	client, err := ethclient.Dial(dsn)
	if err != nil {
		log.Fatal(err)
	}

	return &Ethereum{
		Client: client,
	}
}

func (e *Ethereum) Close() {
	e.Client.Close()
}

func (e *Ethereum) GetCurrentBlock() uint64 {
	blockNumber, err := e.Client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return blockNumber
}

func (e *Ethereum) GetBlockByNumber(blockNumber *big.Int) *types.Block {
	block, err := e.Client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	
	return block
}
