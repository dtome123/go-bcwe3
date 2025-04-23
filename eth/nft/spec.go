package nft

import "github.com/dtome123/go-bcwe3/eth/types"

type NFT interface {
	GetWalletNFTs(account string, contract string) ([]*types.NFT, error)
}
