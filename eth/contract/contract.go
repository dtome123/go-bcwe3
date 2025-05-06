package contract

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dtome123/go-bcwe3/eth/provider"
	"github.com/dtome123/go-bcwe3/eth/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Contract struct {
	provider provider.Provider
}

type Cmd struct {
	address  string
	caller   *bind.BoundContract
	provider provider.Provider
}

// NewContract initializes the Contract module with a given provider.
func NewContract(provider provider.Provider) (Contract, error) {
	return Contract{provider: provider}, nil
}

func (c *Contract) NewCaller(address, abiData string) (*bind.BoundContract, error) {
	address = strings.TrimSpace(address)
	if !common.IsHexAddress(address) {
		return nil, fmt.Errorf("invalid contract address: %s", address)
	}

	parsedABI, err := abi.JSON(strings.NewReader(abiData))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}

	client := c.provider.Client()
	caller := bind.NewBoundContract(common.HexToAddress(address), parsedABI, client, client, client)

	return caller, nil
}

// NewCmd initializes a smart contract interface from address and ABI data.
func (c *Contract) NewCmd(address, abiData string) (*Cmd, error) {

	caller, err := c.NewCaller(address, abiData)
	if err != nil {
		return nil, err
	}

	return &Cmd{
		address:  address,
		caller:   caller,
		provider: c.provider,
	}, nil
}

// Call invokes a view/pure function on the contract.
func (c Cmd) Call(ctx context.Context, method string, params ...any) ([]any, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if method == "" {
		return nil, errors.New("method name cannot be empty")
	}

	var result []any
	err := c.caller.Call(&bind.CallOpts{Context: ctx}, &result, method, params...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Transact sends a state-changing transaction to the contract.
func (e *Cmd) Transact(ctx context.Context, method string, privateKey string, params ...any) (*types.Tx, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if method == "" {
		return nil, errors.New("method cannot be empty")
	}
	if strings.TrimSpace(privateKey) == "" {
		return nil, errors.New("private key is required")
	}

	privateKeyObj, err := crypto.HexToECDSA(strings.TrimPrefix(privateKey, "0x"))
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %v", err)
	}

	chainID, err := e.provider.Client().NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get network ID: %v", err)
	}

	auth := bind.NewKeyedTransactor(privateKeyObj, chainID)

	tx, err := e.caller.Transact(auth, method, params...)
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %v", err)
	}

	return types.WrapTx(tx), nil
}
