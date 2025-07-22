package main

import (
	"context"
	"fmt"

	"github.com/dtome123/go-bcwe3/eth"
)

func main() {

	eth := eth.NewEth("")

	price, _ := eth.GetProvider().SuggestGasPrice(context.Background())

	fmt.Println(price)
}
