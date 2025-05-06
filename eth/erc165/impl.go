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

type ERC165 interface {
	SupportInterface(ctx context.Context, contractAddr string, interfaceIdBytes [4]byte) (bool, error)
}

func New(provider provider.Provider, contract contract.Contract) ERC165 {
	return &impl{
		provider: provider,
		Contract: contract,
	}
}

func (i *impl) SupportInterface(ctx context.Context, contractAddr string, interfaceIdBytes [4]byte) (bool, error) {

	caller, err := i.NewCaller(contractAddr, constants.ERC165ABI)

	if err != nil {
		return false, err
	}

	isMatch, err := contract.CallViewFunction[bool](caller, "supportsInterface", interfaceIdBytes)

	if err != nil {
		return false, err
	}

	return isMatch, nil
}
