package erc20

import (
	"context"
	"math/big"
	"strings"

	"github.com/dtome123/go-bcwe3/eth/constants"
	"github.com/dtome123/go-bcwe3/eth/contract"
	"github.com/dtome123/go-bcwe3/eth/provider"
	"github.com/dtome123/go-bcwe3/eth/types"
	"github.com/ethereum/go-ethereum/common"
)

type impl struct {
	provider provider.Provider
	address  string
	contract contract.Contract
}

func (i *impl) Address() string {
	return i.address
}

func (i *impl) Name() (string, error) {
	result, err := i.contract.CallViewFunction("name")

	if err != nil {
		return "", err
	}

	return result.Index(0).AsString()
}

func (i *impl) Symbol() (string, error) {
	result, err := i.contract.CallViewFunction("symbol")

	if err != nil {
		return "", err
	}

	return result.Index(0).AsString()
}

func (i *impl) Decimals() (uint8, error) {
	result, err := i.contract.CallViewFunction("decimals")

	if err != nil {
		return 0, err
	}

	return result.Index(0).AsUnit8()
}

func (i *impl) TotalSupply() (*big.Int, error) {
	result, err := i.contract.CallViewFunction("totalSupply")

	if err != nil {
		return nil, err
	}

	return result.Index(0).AsBigInt()
}

func (i *impl) BalanceOf(account string) (*big.Int, error) {
	result, err := i.contract.CallViewFunction("balanceOf", common.HexToAddress(account))

	if err != nil {
		return nil, err
	}

	return result.Index(0).AsBigInt()
}

func (i *impl) Transfer(ctx context.Context, toAddress string, amount *big.Int, privateKey string) (*types.Tx, error) {

	tx, err := i.contract.Transact(ctx, "transfer", privateKey, common.HexToAddress(toAddress), amount)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (i *impl) GetInfo(ctx context.Context) (*types.ERC20Token, error) {

	name, err := i.Name()
	if err != nil {
		return nil, err
	}

	symbol, err := i.Symbol()
	if err != nil {
		return nil, err
	}

	decimals, err := i.Decimals()
	if err != nil {
		return nil, err
	}

	totalSupply, err := i.TotalSupply()
	if err != nil {
		return nil, err
	}

	return &types.ERC20Token{
		Name:        name,
		Symbol:      symbol,
		Decimals:    decimals,
		TotalSupply: totalSupply,
		Address:     i.address,
	}, nil
}

func (i *impl) IsPossiblyERC20(ctx context.Context) (bool, error) {

	bytecode, err := i.provider.CodeAt(ctx, i.address, nil)
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

func New(address string, provider provider.Provider) (ERC20, error) {

	contract, err := contract.NewContract(provider, address, constants.ERC20ABI)

	if err != nil {
		return nil, err
	}

	return &impl{
		provider: provider,
		contract: contract,
	}, nil
}
