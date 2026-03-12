package main

import (
	"context"
	"fmt"

	"github.com/dtome123/go-bcwe3/eth"
)

func main() {

	client := eth.NewClient("http://118.69.78.91:8549")

	defer client.Close()

	token, err := client.NewERC20("0xfD9A18b0E43ECEc17DF2eDbbEC2b4936aE07B8db")
	info, err := token.GetInfo(context.Background())

	if err != nil {
		panic(err)
	}

	fmt.Println("name: ", info.Name)
	fmt.Println("symbol: ", info.Symbol)
	fmt.Println("decimals: ", info.Decimals)
	fmt.Println("totalSupply: ", info.TotalSupply)
}
