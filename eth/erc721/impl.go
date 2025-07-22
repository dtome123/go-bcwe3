package erc721

import (
	"context"
	"math/big"
	"sort"

	"github.com/dtome123/go-bcwe3/eth/constants"
	"github.com/dtome123/go-bcwe3/eth/contract"
	"github.com/dtome123/go-bcwe3/eth/erc165"
	"github.com/dtome123/go-bcwe3/eth/provider"
	"github.com/dtome123/go-bcwe3/eth/types"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
) (ERC721, error) {

	erc165, err := erc165.New(address, provider)
	if err != nil {
		return nil, err
	}

	contract, err := contract.NewContract(provider, address, constants.ERC721ABI)
	if err != nil {
		return nil, err
	}

	return &impl{
		provider: provider,
		contract: contract,
		address:  address,
		ERC165:   erc165,
	}, nil
}

func (i *impl) GetOwnerTokens(ctx context.Context) ([]*types.NFTBalance, error) {

	addressHash := common.HexToAddress(i.address)

	eventSignature := []byte("Transfer(address,address,uint256)")
	eventID := crypto.Keccak256Hash(eventSignature)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{addressHash},
		Topics:    [][]common.Hash{{eventID}},
	}

	logs, err := i.provider.FilterLogs(ctx, query)
	if err != nil {
		return nil, err
	}

	sort.Slice(logs, func(i, j int) bool {
		return logs[i].BlockNumber < logs[j].BlockNumber
	})
	owners := make(map[string]string)

	for _, log := range logs {

		if len(log.Topics) != 4 {
			continue
		}

		to := common.HexToAddress(log.Topics[2].Hex())
		tokenId := new(big.Int).SetBytes(log.Topics[3].Bytes()).String()

		if to == (common.Address{}) {
			delete(owners, tokenId)
		} else {
			owners[tokenId] = to.Hex()
		}
	}

	out := make([]*types.NFTBalance, 0, len(owners))
	for tokenId, owner := range owners {
		out = append(out, &types.NFTBalance{
			Token: types.NFT{
				ContractAddress: i.address,
				TokenId:         tokenId,
				Standard:        types.ERC721,
			},
			Owner:   owner,
			Balance: big.NewInt(1),
		})
	}

	return out, nil
}

func (i *impl) IsERC721(ctx context.Context, contractAddr string) (bool, error) {
	// ERC-721 interfaceId is "0x80ac58cd"
	interfaceIdBytes := [4]byte{0x80, 0xac, 0x58, 0xcd}

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
