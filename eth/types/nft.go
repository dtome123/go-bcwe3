package types

type NFTCollection struct {
	ContractAddress string `json:"contract_address"`
	Tokens          []*NFT `json:"tokens"`
}

type NFT struct {
	ContractAddress string `json:"address"`
	TokenId         string `json:"token_id"`
	Uri             string `json:"uri"`
}
