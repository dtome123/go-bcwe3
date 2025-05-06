package erc20

import (
	"context"
	"math/big"

	"github.com/dtome123/go-bcwe3/eth/constants"
	"github.com/dtome123/go-bcwe3/eth/contract"
	"github.com/dtome123/go-bcwe3/eth/provider"
	"github.com/dtome123/go-bcwe3/eth/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type cmd struct {
	address  string
	caller   *bind.BoundContract
	provider provider.Provider
}

type ERC20Cmd interface {
	Address() string
	Name() (string, error)
	Symbol() (string, error)
	Decimals() (uint8, error)
	TotalSupply() (*big.Int, error)
	BalanceOf(address string) (*big.Int, error)
	Transfer(toAddress string, amount *big.Int, privateKey string) (*types.Tx, error)
}

func (i *impl) NewCmd(address string) (ERC20Cmd, error) {

	caller, err := i.NewCaller(address, constants.ERC20ABI)
	if err != nil {
		return nil, err
	}

	return &cmd{
		address:  address,
		caller:   caller,
		provider: i.provider,
	}, nil
}

func (i *cmd) Address() string {
	return i.address
}

func (i *cmd) Name() (string, error) {
	return contract.CallViewFunction[string](i.caller, "name")
}

func (i *cmd) Symbol() (string, error) {
	return contract.CallViewFunction[string](i.caller, "symbol")
}

func (i *cmd) Decimals() (uint8, error) {
	return contract.CallViewFunction[uint8](i.caller, "decimals")
}

func (i *cmd) TotalSupply() (*big.Int, error) {
	return contract.CallViewFunction[*big.Int](i.caller, "totalSupply")
}

func (i *cmd) BalanceOf(address string) (*big.Int, error) {
	return contract.CallViewFunction[*big.Int](i.caller, "balanceOf", common.HexToAddress(address))
}

func (e *cmd) Transfer(toAddress string, amount *big.Int, privateKey string) (*types.Tx, error) {

	privateKeyObj, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	chainID, err := e.provider.Client().NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	auth := bind.NewKeyedTransactor(privateKeyObj, chainID)

	tx, err := e.caller.Transact(auth, "transfer", common.HexToAddress(toAddress), amount)
	if err != nil {
		return nil, err
	}

	return types.WrapTx(tx), nil
}
