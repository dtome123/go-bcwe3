package eth

import (
	"github.com/dtome123/go-bcwe3/eth/contract"
	"github.com/dtome123/go-bcwe3/eth/erc1155"
	"github.com/dtome123/go-bcwe3/eth/erc20"
	"github.com/dtome123/go-bcwe3/eth/erc721"
	"github.com/dtome123/go-bcwe3/eth/listener"
	"github.com/dtome123/go-bcwe3/eth/provider"
)

type Eth struct {
	Provider provider.Provider
	ERC721   erc721.ERC721
	ERC1155  erc1155.ERC1155
	ERC20    erc20.ERC20
	Listener listener.Listener
	Contract contract.Contract
}

func NewEth(dsn string) *Eth {

	provider := provider.NewProvider(dsn)
	contract, _ := contract.NewContract(provider)

	return &Eth{
		Provider: provider,
		Listener: listener.NewListener(provider),
		Contract: contract,

		ERC721:  erc721.New(provider, contract),
		ERC1155: erc1155.New(provider, contract),
		ERC20:   erc20.New(provider, contract),
	}
}
