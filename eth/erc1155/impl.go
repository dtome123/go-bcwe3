package erc1155

import (
	"context"
	"log"
	"math/big"
	"strings"
	"sync"

	"container/heap"

	"github.com/dtome123/go-bcwe3/eth/constants"
	"github.com/dtome123/go-bcwe3/eth/contract"
	"github.com/dtome123/go-bcwe3/eth/erc165"
	"github.com/dtome123/go-bcwe3/eth/provider"
	"github.com/dtome123/go-bcwe3/eth/types"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type impl struct {
	provider provider.Provider
	contract.Contract
	erc165.ERC165
}

func New(
	provider provider.Provider,
	contract contract.Contract,
) ERC1155 {
	erc165 := erc165.New(provider, contract)
	return &impl{
		provider: provider,
		Contract: contract,
		ERC165:   erc165,
	}
}

func (i *impl) GetWalletNFTs(ctx context.Context, account string) ([]*types.NFTCollection, error) {
	targetAddress := common.HexToAddress(account)

	parsedABI, err := abi.JSON(strings.NewReader(constants.ERC1155ABI))
	if err != nil {
		log.Fatal(err)
	}

	singleSig := parsedABI.Events["TransferSingle"].ID
	batchSig := parsedABI.Events["TransferBatch"].ID

	// Fetch logs for "from" and "to" addresses concurrently
	logs, err := i.fetchTransferLogs(singleSig, batchSig, targetAddress)
	if err != nil {
		return nil, err
	}

	if len(logs) == 0 {
		return nil, nil
	}

	// Process logs to determine NFT holdings
	nftHoldings := i.processLogs(singleSig, batchSig, logs, targetAddress)

	if len(nftHoldings) == 0 {
		return nil, nil
	}

	// Query IsERC1155 for each contract in parallel using a worker pool
	isNFTContracts, err := i.checkIsNFTContracts(nftHoldings)
	if err != nil {
		return nil, err
	}

	// Prepare final list of NFTs
	nfts := i.prepareNFTs(nftHoldings, isNFTContracts)
	return nfts, nil
}

func (i *impl) IsERC1155(ctx context.Context, contractAddr string) (bool, error) {
	// ERC-1155 interfaceId is "0xd9b67a26"
	interfaceIdBytes := [4]byte{0xd9, 0xb6, 0x7a, 0x26}

	return i.SupportInterface(ctx, contractAddr, interfaceIdBytes)
}

func (i *impl) GetBalanceOf(ctx context.Context, account string, tokenAddress string) (*big.Int, error) {
	caller, err := i.NewCaller(tokenAddress, constants.ERC1155ABI)

	if err != nil {
		return nil, err
	}

	return contract.CallViewFunction[*big.Int](caller, "balanceOf", common.HexToAddress(account))
}

func (i *impl) GetOwnerOf(ctx context.Context, tokenAddress string, tokenId *big.Int) (string, error) {
	caller, err := i.NewCaller(tokenAddress, constants.ERC1155ABI)

	if err != nil {
		return "", err
	}

	address, err := contract.CallViewFunction[common.Address](caller, "ownerOf", tokenId)

	if err != nil {
		return "", err
	}

	return address.Hex(), nil
}

func (i *impl) GetName(ctx context.Context, tokenAddress string) (string, error) {
	caller, err := i.NewCaller(tokenAddress, constants.ERC1155ABI)

	if err != nil {
		return "", err
	}

	name, err := contract.CallViewFunction[string](caller, "name")

	if err != nil {
		return "", err
	}

	return name, nil
}

func (i *impl) GetSymbol(ctx context.Context, tokenAddress string) (string, error) {
	caller, err := i.NewCaller(tokenAddress, constants.ERC1155ABI)

	if err != nil {
		return "", err
	}

	symbol, err := contract.CallViewFunction[string](caller, "symbol")

	if err != nil {
		return "", err
	}

	return symbol, nil
}

///////////////////////////// private method ////////////////////////////

// fetchTransferLogs fetches logs for both "from" and "to" addresses concurrently
func (i *impl) fetchTransferLogs(singleSig, batchSig common.Hash, targetAddress common.Address) ([]types.Log, error) {

	type result struct {
		logs []types.Log
		err  error
	}

	logsChan := make(chan result, 2)
	startBlock := big.NewInt(0)

	// Query for 'from' address
	go func() {
		query := ethereum.FilterQuery{
			FromBlock: startBlock,
			Topics: [][]common.Hash{
				{singleSig, batchSig},
				nil,
				{common.HexToHash(targetAddress.Hex())},
			},
		}
		logs, err := i.provider.FilterLogs(context.Background(), query)
		logsChan <- result{logs: logs, err: err}
	}()

	// Query for 'to' address
	go func() {
		query := ethereum.FilterQuery{
			FromBlock: startBlock,
			Topics: [][]common.Hash{
				{singleSig, batchSig},
				nil,
				nil,
				{common.HexToHash(targetAddress.Hex())},
			},
		}
		logs, err := i.provider.FilterLogs(context.Background(), query)
		logsChan <- result{logs: logs, err: err}
	}()

	// Collect logs
	var logs []types.Log
	for i := 0; i < 2; i++ {
		res := <-logsChan
		if res.err != nil {
			return nil, res.err
		}
		logs = append(logs, res.logs...)
	}
	close(logsChan)

	return logs, nil
}

// processLogs processes logs and tracks NFT holdings for each contract
func (i *impl) processLogs(singleSig, batchSig common.Hash, logs []types.Log, targetAddress common.Address) map[string]map[string]*big.Int {

	logsHeap := &types.LogHeap{}
	heap.Init(logsHeap)

	// Add logs to heap for sorting
	for _, vLog := range logs {
		heap.Push(logsHeap, vLog)
	}

	type tokenKey struct {
		Contract string
		TokenID  string
	}

	balanceMap := make(map[tokenKey]*big.Int)
	// Process sorted logs
	for logsHeap.Len() > 0 {
		vLog := heap.Pop(logsHeap).(types.Log)

		switch vLog.Topics[0] {
		case singleSig:
			var event struct {
				Id    *big.Int
				Value *big.Int
			}

			from := common.HexToAddress(vLog.Topics[2].Hex())
			to := common.HexToAddress(vLog.Topics[3].Hex())
			key := tokenKey{Contract: vLog.Address, TokenID: event.Id.String()}

			if from == targetAddress {
				if _, ok := balanceMap[key]; !ok {
					balanceMap[key] = big.NewInt(0)
				}
				balanceMap[key].Sub(balanceMap[key], event.Value)
			}
			if to == targetAddress {
				if _, ok := balanceMap[key]; !ok {
					balanceMap[key] = big.NewInt(0)
				}
				balanceMap[key].Add(balanceMap[key], event.Value)
			}

		case batchSig:
			var event struct {
				Ids    []*big.Int
				Values []*big.Int
			}

			from := common.HexToAddress(vLog.Topics[2].Hex())
			to := common.HexToAddress(vLog.Topics[3].Hex())

			for i := 0; i < len(event.Ids); i++ {
				key := tokenKey{Contract: vLog.Address, TokenID: event.Ids[i].String()}
				if _, ok := balanceMap[key]; !ok {
					balanceMap[key] = big.NewInt(0)
				}
				if from == targetAddress {
					balanceMap[key].Sub(balanceMap[key], event.Values[i])
				}
				if to == targetAddress {
					balanceMap[key].Add(balanceMap[key], event.Values[i])
				}
			}
		}
	}

	out := make(map[string]map[string]*big.Int)
	for k, v := range balanceMap {
		if _, ok := out[k.Contract]; !ok {
			out[k.Contract] = make(map[string]*big.Int)
		}
		out[k.Contract][k.TokenID] = v
	}

	return out
}

// checkIsNFTContracts checks if each contract is an ERC-1155 token using parallel workers
func (i *impl) checkIsNFTContracts(nftHoldings map[string]map[string]*big.Int) (
	map[string]bool,
	error,
) {
	workerPool := make(chan struct{}, 10) // Max 10 workers
	isNFTContracts := make(map[string]bool)

	var wg sync.WaitGroup
	for contract := range nftHoldings {
		wg.Add(1)
		go func(contract string) {
			defer wg.Done()
			// Acquire a worker from the pool
			workerPool <- struct{}{}
			defer func() {
				// Release the worker after the task is done
				<-workerPool
			}()

			isNFT, err := i.IsERC1155(context.Background(), contract)
			if err != nil {
				return
			}
			isNFTContracts[contract] = isNFT
		}(contract)
	}

	// Wait for all workers to finish
	wg.Wait()

	return isNFTContracts, nil
}

// prepareNFTs prepares the final list of NFTs based on holdings and ERC-1155 verification
func (i *impl) prepareNFTs(
	nftHoldings map[string]map[string]*big.Int,
	isNFTContracts map[string]bool,
) []*types.NFTCollection {

	collections := make([]*types.NFTCollection, 0)
	for contract, tokens := range nftHoldings {
		if !isNFTContracts[contract] {
			continue
		}

		var nfts []*types.NFT
		for tokenID := range tokens {
			nfts = append(nfts, &types.NFT{
				ContractAddress: contract,
				TokenId:         tokenID,
				Standard:        types.ERC1155,
			})
		}

		collections = append(collections, &types.NFTCollection{
			ContractAddress: contract,
			Tokens:          nfts,
		})
	}

	return collections
}
