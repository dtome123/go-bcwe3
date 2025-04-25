package nft

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/dtome123/go-bcwe3/eth/provider"
	"github.com/dtome123/go-bcwe3/eth/types"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/semaphore"
)

type impl struct {
	provider provider.Provider
}

func NewNFT(provider provider.Provider) NFT {

	return &impl{
		provider: provider,
	}
}

func (n *impl) GetWalletNFTs(account string, contract string) ([]*types.NFT, error) {

	const maxConcurrentCalls = 10

	parsedABI, err := abi.JSON(strings.NewReader(erc721ABI))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParseABI, err)
	}

	accountAddr := common.HexToAddress(account)
	contractAddr := common.HexToAddress(contract)

	balance, err := n.callBalanceOf(parsedABI, accountAddr, contractAddr)
	if err != nil {
		return nil, err
	}

	numTokens := int(balance.Int64())
	tokens := make([]*types.NFT, numTokens)
	var wg sync.WaitGroup
	sem := semaphore.NewWeighted(maxConcurrentCalls)
	errChan := make(chan error, 1)

	for i := range numTokens {
		i := i // capture loop variable
		if err := sem.Acquire(context.Background(), 1); err != nil {
			return nil, fmt.Errorf("failed to acquire semaphore: %w", err)
		}
		wg.Add(1)

		go func() {
			defer wg.Done()
			defer sem.Release(1)

			tokenID, err := n.callTokenOfOwnerByIndex(parsedABI, accountAddr, contractAddr, big.NewInt(int64(i)))
			if err != nil {
				select {
				case errChan <- fmt.Errorf("index %d: %w", i, err):
				default:
				}
				return
			}
			tokens[i] = &types.NFT{
				TokenId: tokenID,
				Address: contract,
			}
		}()
	}

	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return nil, err
	}

	return tokens, nil
}

/////////////////////////////// PRIVATE ////////////////////////////////

func (n *impl) callBalanceOf(parsedABI abi.ABI, owner, contract common.Address) (*big.Int, error) {
	data, err := parsedABI.Pack("balanceOf", owner)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPackBalanceOf, err)
	}

	res, err := n.provider.CallContract(context.Background(), ethereum.CallMsg{
		To: &contract, Data: data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCallBalanceOf, err)
	}
	if len(res) == 0 {
		return nil, ErrEmptyBalanceOfResponse
	}

	var balance = new(big.Int)
	if err := parsedABI.UnpackIntoInterface(&balance, "balanceOf", res); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnpackBalanceOf, err)
	}
	return balance, nil
}

func (n *impl) callTokenOfOwnerByIndex(parsedABI abi.ABI, owner, contract common.Address, index *big.Int) (*big.Int, error) {
	data, err := parsedABI.Pack("tokenOfOwnerByIndex", owner, index)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrPackTokenOfOwnerByIndex, err)
	}

	res, err := n.provider.CallContract(context.Background(), ethereum.CallMsg{
		To: &contract, Data: data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCallTokenOfOwnerByIndex, err)
	}
	if len(res) == 0 {
		return nil, ErrEmptyTokenOfOwnerResponse
	}

	var tokenID = new(big.Int)
	if err := parsedABI.UnpackIntoInterface(&tokenID, "tokenOfOwnerByIndex", res); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnpackTokenOfOwnerByIndex, err)
	}
	return tokenID, nil
}
