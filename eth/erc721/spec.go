package erc721

import (
	"context"
	"math/big"

	"github.com/dtome123/go-bcwe3/eth/types"
)

type ERC721 interface {
	IsERC721(ctx context.Context, contractAddr string) (bool, error)
	GetWalletNFTs(ctx context.Context, account string) ([]*types.NFTCollection, error)
	GetOwnerTokens(ctx context.Context, tokenAddress string) ([]*types.NFTBalance, error)
	GetBalanceOf(ctx context.Context, account string, tokenAddress string) (*big.Int, error)
	GetOwnerOf(ctx context.Context, tokenAddress string, tokenId *big.Int) (string, error)
	GetName(ctx context.Context, tokenAddress string) (string, error)
	GetSymbol(ctx context.Context, tokenAddress string) (string, error)
}
