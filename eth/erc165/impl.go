package erc165

import (
	"context"

	"github.com/dtome123/go-bcwe3/eth/constants"
	"github.com/dtome123/go-bcwe3/eth/contract"
	"github.com/dtome123/go-bcwe3/eth/provider"
)

type impl struct {
	provider provider.Provider
	contract.Contract
}

func New(address string, provider provider.Provider) (ERC165, error) {

	contract, err := contract.NewContract(provider, address, constants.ERC165ABI)

	if err != nil {
		return nil, err
	}

	return &impl{
		provider: provider,
		Contract: contract,
	}, nil
}

func (i *impl) SupportInterface(ctx context.Context, contractAddr string, interfaceIdBytes [4]byte) (bool, error) {

	result, err := i.Contract.Call(context.Background(), "supportsInterface", interfaceIdBytes)

	if err != nil {
		return false, err
	}

	return result.Index(0).AsBool()
}
