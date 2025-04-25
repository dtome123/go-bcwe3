package nft

import (
	"context"
	"log"
	"math/big"
	"strings"
	"sync"

	"github.com/dtome123/go-bcwe3/eth/provider"
	"github.com/dtome123/go-bcwe3/eth/types"

	"container/heap"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	goethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type impl struct {
	provider provider.Provider
}

func NewNFT(provider provider.Provider) NFT {

	return &impl{
		provider: provider,
	}
}

func (i *impl) GetWalletNFTs(account string) ([]*types.NFTCollection, error) {
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

func (i *impl) IsERC721(contractAddr string) (bool, error) {
	// ERC-721 interfaceId is "0x80ac58cd"
	interfaceIdBytes := [4]byte{0x80, 0xac, 0x58, 0xcd}

	return i.supportInterface(contractAddr, interfaceIdBytes)
}

func (i *impl) IsERC1155(contractAddr string) (bool, error) {
	// ERC-1155 interfaceId is "0xd9b67a26"
	interfaceIdBytes := [4]byte{0xd9, 0xb6, 0x7a, 0x26}

	return i.supportInterface(contractAddr, interfaceIdBytes)
}

func (i *impl) IsNFTToken(contractAddr string) (bool, error) {

	var isNFT bool
	var err error

	isNFT, err = i.IsERC721(contractAddr)
	if err != nil {
		return false, err
	}
	if isNFT {
		return true, nil
	}

	isNFT, err = i.IsERC1155(contractAddr)
	if err != nil {
		return false, err
	}
	if isNFT {
		return true, nil
	}

	return false, nil
}

///////////////////////////// private method ////////////////////////////

func (i *impl) supportInterface(contractAddr string, interfaceIdBytes [4]byte) (bool, error) {
	// Convert the contract address from string to Address
	contractAddress := common.HexToAddress(contractAddr)

	// Parse the ERC-165 ABI
	contractABI, err := abi.JSON(strings.NewReader(erc165ABI))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	// Pack the "supportsInterface" method with the interfaceIdBytes as argument
	callData, err := contractABI.Pack("supportsInterface", interfaceIdBytes)
	if err != nil {
		log.Fatalf("Failed to pack data: %v", err)
	}

	// Send the request to the contract and get the result
	result, err := i.provider.CallContract(context.Background(), ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}, nil)
	if err != nil {
		log.Fatalf("Failed to call contract: %v", err)
	}

	// Decode the result (true/false)
	var supports bool
	err = contractABI.UnpackIntoInterface(&supports, "supportsInterface", result)
	if err != nil {
		log.Fatalf("Failed to unpack result: %v", err)
	}

	// Return whether the contract supports ERC-721
	return supports, nil
}

// fetchTransferLogs fetches logs for both "from" and "to" addresses concurrently
func (i *impl) fetchTransferLogs(eventID common.Hash, targetAddress common.Address) ([]goethTypes.Log, error) {

	type result struct {
		logs []goethTypes.Log
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
				{common.BytesToHash(targetAddress.Bytes())},
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
				{common.BytesToHash(targetAddress.Bytes())},
			},
		}
		logs, err := i.provider.FilterLogs(context.Background(), query)
		logsChan <- result{logs: logs, err: err}
	}()

	// Collect logs
	var logs []goethTypes.Log
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
func (i *impl) processLogs(logs []goethTypes.Log, targetAddress common.Address) map[string]map[string]bool {
	nftHoldings := make(map[string]map[string]bool)

	logsHeap := &LogHeap{}
	heap.Init(logsHeap)

	// Add logs to heap for sorting
	for _, vLog := range logs {
		heap.Push(logsHeap, vLog)
	}

	// Process sorted logs
	for logsHeap.Len() > 0 {
		vLog := heap.Pop(logsHeap).(goethTypes.Log)

		if len(vLog.Topics) < 4 {
			continue
		}

		from := common.HexToAddress(vLog.Topics[1].Hex())
		to := common.HexToAddress(vLog.Topics[2].Hex())
		tokenID := new(big.Int).SetBytes(vLog.Topics[3].Bytes()).String()
		contract := vLog.Address.Hex()

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
	}

	return nftHoldings
}

// checkIsNFTContracts checks if each contract is an ERC-721 token using parallel workers
func (i *impl) checkIsNFTContracts(nftHoldings map[string]map[string]bool) (map[string]bool, error) {
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

			isNFT, err := i.IsNFTToken(contract)
			if err != nil {
				// Handle error or log it
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
func (i *impl) prepareNFTs(nftHoldings map[string]map[string]bool, isNFTContracts map[string]bool) []*types.NFTCollection {

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
			})
		}

		collections = append(collections, &types.NFTCollection{
			ContractAddress: contract,
			Tokens:          nfts,
		})
	}

	return collections
}
