package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dtome123/go-bcwe3/eth"
	"github.com/dtome123/go-bcwe3/eth/types"
)

func main() {
	// account := "0x7556989c2A60E60F0c66A2b9D77079BC9F189037"
	// contract := "0x8dbB1977011A586c5F3a58AaC9A07e8CF9eBc0Fd"

	eth := eth.NewEth("wss://sepolia.infura.io/ws/v3/da05d3dc31244bd483a28d746233d32f")

	errChan := make(chan error, 1)

	handleBlock := func(block *types.Block) {
		fmt.Println("Processing block:", block.Number)
	}

	go eth.Listener.ListenEventBlock(handleBlock, errChan)

	go func() {
		for err := range errChan {
			log.Printf("Error: %v", err)
		}
	}()

	// Giả lập chạy ứng dụng trong một khoảng thời gian
	time.Sleep(2 * time.Minute)

}
