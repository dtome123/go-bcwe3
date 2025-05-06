package types

import "math/big"

type ERC20Token struct {
	Address     string   `json:"address"`
	Symbol      string   `json:"symbol"`
	Name        string   `json:"name"`
	Decimals    uint8    `json:"decimals"`
	TotalSupply *big.Int `json:"total_supply"`
}

type ERC20Balance struct {
	Token   ERC20Token
	Balance *big.Int
}
