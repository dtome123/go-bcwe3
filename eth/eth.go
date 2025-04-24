package eth

import (
	"github.com/dtome123/go-bcwe3/eth/client"
	"github.com/dtome123/go-bcwe3/eth/listener"
	"github.com/dtome123/go-bcwe3/eth/nft"
)

type Eth struct {
	Client   client.Client
	NFT      nft.NFT
	Listener listener.Listener
}

func NewEth(dsn string) *Eth {

	client := client.NewClient(dsn)

	return &Eth{
		Client:   client,
		NFT:      nft.NewNFT(client),
		Listener: listener.NewListener(client),
	}
}
