package erc20

import (
	"context"

	"github.com/dtome123/go-bcwe3/eth/types"
)

type ERC20 interface {
	GetInfo(ctx context.Context, contractAddress string) (*types.ERC20Token, error)
	IsPossiblyERC20(ctx context.Context, address string) (bool, error)
	NewCmd(address string) (ERC20Cmd, error)
}
