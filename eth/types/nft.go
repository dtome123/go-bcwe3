package types

import "math/big"

type NFT struct {
	Address string   `json:"address"`
	TokenId *big.Int `json:"token_id"`
}
