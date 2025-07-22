package erc1155

import (
	"context"
	"math/big"

	"github.com/dtome123/go-bcwe3/eth/constants"
	"github.com/dtome123/go-bcwe3/eth/contract"
	"github.com/dtome123/go-bcwe3/eth/erc165"
	"github.com/dtome123/go-bcwe3/eth/provider"

	"github.com/ethereum/go-ethereum/common"
)

type impl struct {
	provider provider.Provider
	address  string
	contract contract.Contract
	erc165.ERC165
}

func New(
	address string,
	provider provider.Provider,
) (ERC1155, error) {
	erc165, err := erc165.New(address, provider)
	if err != nil {
		return nil, err
	}

	contract, err := contract.NewContract(provider, address, constants.ERC1155ABI)
	if err != nil {
		return nil, err
	}
	return &impl{
		provider: provider,
		address:  address,
		contract: contract,
		ERC165:   erc165,
	}, nil
}

func (i *impl) IsERC1155(ctx context.Context, contractAddr string) (bool, error) {
	// ERC-1155 interfaceId is "0xd9b67a26"
	interfaceIdBytes := [4]byte{0xd9, 0xb6, 0x7a, 0x26}

	return i.SupportInterface(ctx, contractAddr, interfaceIdBytes)
}

func (i *impl) GetBalanceOf(ctx context.Context, account string) (*big.Int, error) {

	result, err := i.contract.CallViewFunction("balanceOf", common.HexToAddress(account))

	if err != nil {
		return nil, err
	}

	return result.Index(0).AsBigInt()
}

func (i *impl) GetOwnerOf(ctx context.Context, tokenId *big.Int) (string, error) {

	result, err := i.contract.CallViewFunction("ownerOf", tokenId)

	if err != nil {
		return "", err
	}

	return result.Index(0).AsString()
}

func (i *impl) GetName(ctx context.Context) (string, error) {

	result, err := i.contract.CallViewFunction("name")

	if err != nil {
		return "", err
	}

	return result.Index(0).AsString()
}

func (i *impl) GetSymbol(ctx context.Context) (string, error) {

	result, err := i.contract.CallViewFunction("symbol")

	if err != nil {
		return "", err
	}

	return result.Index(0).AsString()
}
