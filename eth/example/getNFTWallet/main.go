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
	defer eth.Close()

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

	erc, err := eth.NewERC721(tokenAddress)
	balance, err := erc.GetBalanceOf(context.Background(), account)

	if err != nil {
		panic(err)
	}

	fmt.Println("balance: ", balance)

	ownerOf0, err := erc.GetOwnerOf(context.Background(), big.NewInt(0))

	if err != nil {
		panic(err)
	}

	fmt.Println("owner token 0: ", ownerOf0)

	name, err := erc.GetName(context.Background())

	if err != nil {
		panic(err)
	}

	fmt.Println("name: ", name)

	symbol, err := erc.GetSymbol(context.Background())

	if err != nil {
		panic(err)
	}

	fmt.Println("symbol: ", symbol)
}
