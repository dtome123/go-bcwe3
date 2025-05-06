package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/dtome123/go-bcwe3/eth"
	"github.com/ethereum/go-ethereum/common"
)

// Định nghĩa sự kiện Transfer của ERC721
type TransferEvent struct {
	From    common.Address
	To      common.Address
	TokenID *big.Int
}

func main() {
	// URL của Ethereum node (Infura hoặc địa chỉ node của bạn)
	infuraURL := "wss://sepolia.infura.io/ws/v3/da05d3dc31244bd483a28d746233d32f"
	contractAddress := "0x975Bda1d9287433B868967685bC79637A08EfBEc"

	eth := eth.NewEth(infuraURL)
	balance, _ := eth.ERC721.GetOwnerTokens(context.Background(), contractAddress)

	for _, nft := range balance {
		fmt.Println(nft.Owner, ":", nft.Token.TokenId)
	}
}
