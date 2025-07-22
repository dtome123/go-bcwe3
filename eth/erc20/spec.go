package erc20

import (
	"context"
	"math/big"

	"github.com/dtome123/go-bcwe3/eth/contract"
	"github.com/dtome123/go-bcwe3/eth/types"
)

type ERC20 interface {
	contract.Contract
	GetInfo(ctx context.Context) (*types.ERC20Token, error)
	IsPossiblyERC20(ctx context.Context) (bool, error)
	Address() string
	Name() (string, error)
	Symbol() (string, error)
	Decimals() (uint8, error)
	TotalSupply() (*big.Int, error)
	BalanceOf(account string) (*big.Int, error)
}
