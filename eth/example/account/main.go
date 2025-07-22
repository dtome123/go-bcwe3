package main

import (
	"github.com/dtome123/go-bcwe3/eth"
)

func main() {

	eth := eth.NewEth("wss://sepolia.infura.io/ws/v3/da05d3dc31244bd483a28d746233d32f")

	defer eth.Provider().Close()

}
