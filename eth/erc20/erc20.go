package erc20

import (
	"context"
	"strings"

	"github.com/dtome123/go-bcwe3/eth/constants"
	"github.com/dtome123/go-bcwe3/eth/contract"
	"github.com/dtome123/go-bcwe3/eth/provider"
	"github.com/dtome123/go-bcwe3/eth/types"
	"github.com/ethereum/go-ethereum/common"
)

type impl struct {
	provider provider.Provider
	contract.Contract
}

func New(provider provider.Provider, contract contract.Contract) ERC20 {

	return &impl{
		provider: provider,
		Contract: contract,
	}
}

func (e *impl) GetInfo(ctx context.Context, contractAddress string) (*types.ERC20Token, error) {

	cmd, _ := e.NewCmd(contractAddress)

	name, err := cmd.Name()
	if err != nil {
		return nil, err
	}

	symbol, err := cmd.Symbol()
	if err != nil {
		return nil, err
	}

	decimals, err := cmd.Decimals()
	if err != nil {
		return nil, err
	}

	totalSupply, err := cmd.TotalSupply()
	if err != nil {
		return nil, err
	}

	return &types.ERC20Token{
		Name:        name,
		Symbol:      symbol,
		Decimals:    decimals,
		TotalSupply: totalSupply,
		Address:     contractAddress,
	}, nil
}

func (i *impl) IsPossiblyERC20(ctx context.Context, address string) (bool, error) {

	bytecode, err := i.provider.CodeAt(ctx, address, nil)
	if err != nil {
		return false, err
	}
	if len(bytecode) == 0 {
		return false, nil
	}

	hexCode := strings.ToLower(common.Bytes2Hex(bytecode))
	matched := []string{}

	for method, selector := range constants.ERC20Selectors {
		if strings.Contains(hexCode, selector) {
			matched = append(matched, method)
		}
	}

	return len(matched) == len(constants.ERC20Selectors), nil
}
