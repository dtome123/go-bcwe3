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

type implContract struct {
	address       string
	provider      provider.Provider
	boundContract *bind.BoundContract
}

// NewContract initializes the implContract module with a given provider.
func NewContract(provider provider.Provider, address, abiData string) (Contract, error) {

	address = strings.TrimSpace(address)
	if !common.IsHexAddress(address) {
		return nil, fmt.Errorf("invalid implContract address: %s", address)
	}

	parsedABI, err := abi.JSON(strings.NewReader(abiData))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %v", err)
	}

	client := provider.Client()
	contract := bind.NewBoundContract(common.HexToAddress(address), parsedABI, client, client, client)

	return &implContract{
		provider:      provider,
		address:       address,
		boundContract: contract,
	}, nil
}

func (c *implContract) Call(ctx context.Context, method string, params ...any) (ContractResults, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if method == "" {
		return nil, errors.New("method name cannot be empty")
	}

	var result []any
	err := c.boundContract.Call(&bind.CallOpts{Context: ctx}, &result, method, params...)
	if err != nil {
		return nil, err
	}

	contractResults := make(ContractResults, len(result))
	for i, value := range result {
		contractResults[i] = ContractResult{Value: value}
	}

	return contractResults, nil
}

// Transact sends a state-changing transaction to the contract.
func (c *implContract) Transact(ctx context.Context, method string, privateKey string, params ...any) (*types.Tx, error) {
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

	chainID, err := c.provider.Client().NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get network ID: %v", err)
	}

	auth := bind.NewKeyedTransactor(privateKeyObj, chainID)

	tx, err := c.boundContract.Transact(auth, method, params...)
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %v", err)
	}

	return types.WrapTx(tx), nil
}
