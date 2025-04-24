package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/dtome123/go-bcwe3/eth"
)

func main() {
	eth := eth.NewEth("https://sepolia.infura.io/v3/da05d3dc31244bd483a28d746233d32f")

	// blockNumber, _ := eth.Client.BlockNumber(context.Background())
	// fmt.Println("Current block number:", blockNumber)
	ctx := context.Background()
	block, _ := eth.Client.BlockByNumber(ctx, big.NewInt(8184711))

	// tx1 := block.Transactions[0]

	// fmt.Println("tx hash: ", tx1.Hash)
	// fmt.Println("tx gas: ", tx1.Gas)
	// fmt.Println("tx gas price: ", tx1.GasPrice)
	// fee, _ := eth.Client.CalculateTxFee(tx1)
	// fmt.Println("tx fee: ", fee)

	// completeTx, _, err := eth.Client.TransactionByHash(ctx, tx1.Hash)
	// if err != nil {
	// 	panic(err)
	// }

	if block == nil {
		panic("block is nil")
	}

	isFinished, _ := eth.Client.IsBlockFinalized(context.Background(), block.Number)
	fmt.Println("isFinished: ", isFinished)

	fmt.Println("tx len: ", len(block.Transactions))
	fmt.Println("with len", len(block.Withdrawals))
	fmt.Println("block hash:", block.Hash)

	block.Transactions = nil
	block.Withdrawals = nil
	jsonCompleteTx, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println("data: ", string(jsonCompleteTx))
}
