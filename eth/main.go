package eth

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Ethereum struct {
	Client *ethclient.Client
}

func NewEthereum() *Ethereum {
	client, err := ethclient.Dial("http://localhost:8545")
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

func (e *Ethereum) GetBlockByNumber(blockNumber *big.Int) {
	block, err := e.Client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(block)
}
