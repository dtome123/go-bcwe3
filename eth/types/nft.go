package types

import "math/big"

type NFTStandard int32

const (
	ERC721  NFTStandard = 1
	ERC1155 NFTStandard = 2
)

type NFTCollection struct {
	ContractAddress string `json:"contract_address"`
	Tokens          []*NFT `json:"tokens"`
}

type NFT struct {
	ContractAddress string      `json:"address"`
	TokenId         string      `json:"token_id"`
	Standard        NFTStandard `json:"standard"`
}

type NFTBalance struct {
	Token   NFT      `json:"token"`
	Balance *big.Int `json:"balance"`
	Owner   string   `json:"owner"`
}
