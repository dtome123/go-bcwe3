package main

import (
	"fmt"

	"github.com/dtome123/go-bcwe3/eth"
)

func main() {
	account := "0x7556989c2A60E60F0c66A2b9D77079BC9F189037"
	contract := "0x8dbB1977011A586c5F3a58AaC9A07e8CF9eBc0Fd"

	eth := eth.NewEth("https://sepolia.infura.io/v3/da05d3dc31244bd483a28d746233d32f")

	nfts, err := eth.NFT.GetWalletNFTs(account, contract)

	if err != nil {
		panic(err)
	}

	for _, nft := range nfts {
		fmt.Println("token id:", nft.TokenId)
	}

}
