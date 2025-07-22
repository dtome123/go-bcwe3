package eth

import (
	"github.com/dtome123/go-bcwe3/eth/contract"
	"github.com/dtome123/go-bcwe3/eth/erc1155"
	"github.com/dtome123/go-bcwe3/eth/erc20"
	"github.com/dtome123/go-bcwe3/eth/erc721"
	"github.com/dtome123/go-bcwe3/eth/provider"
)

type impl struct {
	provider provider.Provider
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

func (eth *impl) GetProvider() provider.Provider {
	return eth.provider
}

func (eth *impl) NewERC721(address string) (erc721.ERC721, error) {
	return erc721.New(address, eth.provider)
}
func (eth *impl) NewERC1155(address string) (erc1155.ERC1155, error) {
	return erc1155.New(address, eth.provider)
}

func (eth *impl) NewERC20(address string) (erc20.ERC20, error) {
	return erc20.New(address, eth.provider)
}

func (eth *impl) NewContract(address string, abiData string) (contract.Contract, error) {
	return contract.NewContract(eth.provider, address, abiData)
}
