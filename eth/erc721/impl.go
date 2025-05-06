package erc721

import (
	"context"
	"math/big"
	"sort"
	"sync"

	"container/heap"

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
	contract.Contract
	erc165.ERC165
}

func New(
	provider provider.Provider,
	contract contract.Contract,
) ERC721 {

	erc165 := erc165.New(provider, contract)
	return &impl{
		provider: provider,
		Contract: contract,
		ERC165:   erc165,
	}
}

func (i *impl) GetWalletNFTs(ctx context.Context, account string) ([]*types.NFTCollection, error) {
	eventID := crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
	targetAddress := common.HexToAddress(account)

	// Fetch logs for "from" and "to" addresses concurrently
	logs, err := i.fetchTransferLogs(eventID, targetAddress)
	if err != nil {
		return nil, err
	}

	if len(logs) == 0 {
		return nil, nil
	}

	// Process logs to determine NFT holdings
	nftHoldings := i.processLogs(logs, targetAddress)

	if len(nftHoldings) == 0 {
		return nil, nil
	}

	// Query IsERC721 for each contract in parallel using a worker pool
	isNFTContracts, err := i.checkIsNFTContracts(nftHoldings)
	if err != nil {
		return nil, err
	}

	// Prepare final list of NFTs
	nfts := i.prepareNFTs(nftHoldings, isNFTContracts)
	return nfts, nil
}

func (i *impl) GetOwnerTokens(ctx context.Context, tokenAddress string) ([]*types.NFTBalance, error) {

	addressHash := common.HexToAddress(tokenAddress)

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
				ContractAddress: tokenAddress,
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

func (i *impl) GetBalanceOf(ctx context.Context, account string, tokenAddress string) (*big.Int, error) {
	caller, err := i.NewCaller(tokenAddress, constants.ERC721ABI)

	if err != nil {
		return nil, err
	}

	return contract.CallViewFunction[*big.Int](caller, "balanceOf", common.HexToAddress(account))
}

func (i *impl) GetOwnerOf(ctx context.Context, tokenAddress string, tokenId *big.Int) (string, error) {
	caller, err := i.NewCaller(tokenAddress, constants.ERC721ABI)

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
	caller, err := i.NewCaller(tokenAddress, constants.ERC721ABI)

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
	caller, err := i.NewCaller(tokenAddress, constants.ERC721ABI)

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
func (i *impl) fetchTransferLogs(eventID common.Hash, targetAddress common.Address) ([]types.Log, error) {

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
				{eventID},
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
				{eventID},
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
func (i *impl) processLogs(logs []types.Log, targetAddress common.Address) map[string]map[string]bool {
	nftHoldings := make(map[string]map[string]bool)

	logsHeap := &types.LogHeap{}
	heap.Init(logsHeap)

	// Add logs to heap for sorting
	for _, vLog := range logs {
		heap.Push(logsHeap, vLog)
	}

	// Process sorted logs
	for logsHeap.Len() > 0 {
		vLog := heap.Pop(logsHeap).(types.Log)

		if len(vLog.Topics) < 4 {
			continue
		}

		from := common.HexToAddress(vLog.Topics[1].Hex())
		to := common.HexToAddress(vLog.Topics[2].Hex())
		tokenID := new(big.Int).SetBytes(vLog.Topics[3].Bytes()).String()
		contract := vLog.Address

		// Initialize contract entry in nftHoldings map if not exists
		if _, ok := nftHoldings[contract]; !ok {
			nftHoldings[contract] = make(map[string]bool)
		}

		// Update NFT holdings based on transfer direction
		if to == targetAddress {
			nftHoldings[contract][tokenID] = true
		}
		if from == targetAddress {
			delete(nftHoldings[contract], tokenID)
		}

		if len(nftHoldings[contract]) == 0 {
			delete(nftHoldings, contract)
		}
	}

	return nftHoldings
}

// checkIsNFTContracts checks if each contract is an ERC-721 token using parallel workers
func (i *impl) checkIsNFTContracts(nftHoldings map[string]map[string]bool) (
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

			isNFT, err := i.IsERC721(context.Background(), contract)
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

// prepareNFTs prepares the final list of NFTs based on holdings and ERC-721 verification
func (i *impl) prepareNFTs(
	nftHoldings map[string]map[string]bool,
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
				Standard:        types.ERC721,
			})
		}

		collections = append(collections, &types.NFTCollection{
			ContractAddress: contract,
			Tokens:          nfts,
		})
	}

	return collections
}
