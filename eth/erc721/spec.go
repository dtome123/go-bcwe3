package erc721

import (
	"context"
	"math/big"

	"github.com/dtome123/go-bcwe3/eth/contract"
	"github.com/dtome123/go-bcwe3/eth/types"
)

type ERC721 interface {
	contract.Contract
	GetOwnerTokens(ctx context.Context) ([]*types.NFTBalance, error)
	IsERC721(ctx context.Context, contractAddr string) (bool, error)
	GetBalanceOf(ctx context.Context, account string) (*big.Int, error)
	GetOwnerOf(ctx context.Context, tokenId *big.Int) (string, error)
	GetName(ctx context.Context) (string, error)
	GetSymbol(ctx context.Context) (string, error)
}
