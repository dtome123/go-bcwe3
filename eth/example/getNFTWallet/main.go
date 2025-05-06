package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/dtome123/go-bcwe3/eth"
	"github.com/ethereum/go-ethereum/common"
)

// TransferInfo holds the basic Transfer event data
type TransferInfo struct {
	From    common.Address
	To      common.Address
	TokenID string
}

func main() {

	eth := eth.NewEth("wss://sepolia.infura.io/ws/v3/da05d3dc31244bd483a28d746233d32f")
	account := "0x7556989c2A60E60F0c66A2b9D77079BC9F189037"
	tokenAddress := "0x975Bda1d9287433B868967685bC79637A08EfBEc"
	defer eth.Provider.Close()

	// start := time.Now()
	// collections, err := eth.ERC721.GetWalletNFTs(context.Background(), "0x7556989c2A60E60F0c66A2b9D77079BC9F189037")

	// if err != nil {
	// 	panic(err)
	// }

	// elapsed := time.Since(start)
	// fmt.Printf("Time taken: %s\n", elapsed)

	// for _, collection := range collections {
	// 	fmt.Println("collection address:", collection.ContractAddress)

	// 	for _, nft := range collection.Tokens {
	// 		fmt.Println("token id:", nft.TokenId)
	// 	}
	// }

	balance, err := eth.ERC721.GetBalanceOf(context.Background(), account, tokenAddress)

	if err != nil {
		panic(err)
	}

	fmt.Println("balance: ", balance)

	ownerOf0, err := eth.ERC721.GetOwnerOf(context.Background(), tokenAddress, big.NewInt(0))

	if err != nil {
		panic(err)
	}

	fmt.Println("owner token 0: ", ownerOf0)


	name, err := eth.ERC721.GetName(context.Background(), tokenAddress)

	if err != nil {
		panic(err)
	}

	fmt.Println("name: ", name)


	symbol, err := eth.ERC721.GetSymbol(context.Background(), tokenAddress)

	if err != nil {
		panic(err)
	}

	fmt.Println("symbol: ", symbol)
}
