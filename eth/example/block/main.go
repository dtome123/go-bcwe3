package main

import (
	"context"
	"fmt"

	"github.com/dtome123/go-bcwe3/eth"
)

func main() {

	client := eth.NewClient("ws://118.69.78.91:8586")
	n, err := client.GetProvider().BlockNumber(context.Background())

	if err != nil {
		panic(err)
	}

	fmt.Println("block num: ", n)
}
