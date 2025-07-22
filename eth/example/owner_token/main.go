package main

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// Định nghĩa sự kiện Transfer của ERC721
type TransferEvent struct {
	From    common.Address
	To      common.Address
	TokenID *big.Int
}

func main() {

}
