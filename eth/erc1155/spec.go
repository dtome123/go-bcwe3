package erc1155

import (
	"context"
	"math/big"
)

type ERC1155 interface {
	IsERC1155(ctx context.Context, contractAddr string) (bool, error)
	GetBalanceOf(ctx context.Context, account string) (*big.Int, error)
	GetOwnerOf(ctx context.Context, tokenId *big.Int) (string, error)
	GetName(ctx context.Context) (string, error)
	GetSymbol(ctx context.Context) (string, error)
}
