package main

import (
	"context"
	"fmt"

	"github.com/dtome123/go-bcwe3/eth"
)

func main() {

	eth := eth.NewEth("wss://sepolia.infura.io/ws/v3/da05d3dc31244bd483a28d746233d32f")

	defer eth.Close()

	res, err := eth.GetProvider().BalanceAt(context.Background(), "0x7556989c2A60E60F0c66A2b9D77079BC9F189037", nil)

	if err != nil {
		panic(err)
	}

	fmt.Println("balance: ", res)
}
