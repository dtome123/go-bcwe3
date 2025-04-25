package nft

import "github.com/dtome123/go-bcwe3/eth/types"

type NFT interface {
	GetWalletNFTs(account string) ([]*types.NFTCollection, error)
	IsERC721(contractAddr string) (bool, error)
	IsERC1155(contractAddr string) (bool, error)
}
