package eth

import (
	"github.com/dtome123/go-bcwe3/eth/erc1155"
	"github.com/dtome123/go-bcwe3/eth/erc20"
	"github.com/dtome123/go-bcwe3/eth/erc721"
	"github.com/dtome123/go-bcwe3/eth/listener"
	"github.com/dtome123/go-bcwe3/eth/provider"
)

type impl struct {
	provider provider.Provider
}

type Eth interface {
	Close()
	Provider() provider.Provider
	ERC721(address string) (erc721.ERC721, error)
	ERC1155(address string) (erc1155.ERC1155, error)
	ERC20(address string) (erc20.ERC20, error)
	Listener() listener.Listener
}

func NewEth(dsn string) Eth {

	provider := provider.NewProvider(dsn)

	return &impl{
		provider: provider,
	}
}

func (eth *impl) Close() {
	eth.provider.Close()
}

func (eth *impl) Provider() provider.Provider {
	return eth.provider
}

func (eth *impl) Listener() listener.Listener {
	return listener.NewListener(eth.provider)
}

func (eth *impl) ERC721(address string) (erc721.ERC721, error) {
	return erc721.New(address, eth.provider)
}
func (eth *impl) ERC1155(address string) (erc1155.ERC1155, error) {
	return erc1155.New(address, eth.provider)
}

func (eth *impl) ERC20(address string) (erc20.ERC20, error) {
	return erc20.New(address, eth.provider)
}
