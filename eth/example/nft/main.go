package main

import (
	"context"
	"fmt"

	"github.com/dtome123/go-bcwe3/eth"
)

func main() {
	account := "0xa84e540D1eb5458DFC2bC25760bD64fbECb8e345"

	eth := eth.NewEth("https://sepolia.infura.io/v3/da05d3dc31244bd483a28d746233d32f")

	collections, err := eth.ERC721.GetWalletNFTs(context.Background(), account)

	if err != nil {
		panic(err)
	}

	for _, collection := range collections {
		for _, nft := range collection.Tokens {
			fmt.Println(nft.ContractAddress, ": ", nft.TokenId)
		}
	}

}
