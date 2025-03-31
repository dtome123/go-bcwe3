package interfaces

type NFT interface {
	GetWalletNFTs(address string)
	GetMultipleNFTs(tokenAddress string, tokenIds []string)
}
