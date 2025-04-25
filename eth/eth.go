package eth

import (
	"github.com/dtome123/go-bcwe3/eth/listener"
	"github.com/dtome123/go-bcwe3/eth/nft"
	"github.com/dtome123/go-bcwe3/eth/provider"
)

type Eth struct {
	Provider provider.Provider
	NFT      nft.NFT
	Listener listener.Listener
}

func NewEth(dsn string) *Eth {

	provider := provider.NewProvider(dsn)

	return &Eth{
		Provider: provider,
		NFT:      nft.NewNFT(provider),
		Listener: listener.NewListener(provider),
	}
}
