package main

import (
	"fmt"
	"time"

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

	defer eth.Provider.Close()

	start := time.Now()
	collections, err := eth.NFT.GetWalletNFTs("0xa84e540D1eb5458DFC2bC25760bD64fbECb8e345")

	if err != nil {
		panic(err)
	}

	elapsed := time.Since(start)
	fmt.Printf("Time taken: %s\n", elapsed)

	for _, collection := range collections {
		fmt.Println("collection address:", collection.ContractAddress)

		for _, nft := range collection.Tokens {
			fmt.Println("token id:", nft.TokenId)
		}
	}

}
