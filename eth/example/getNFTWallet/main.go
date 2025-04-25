package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/dtome123/go-bcwe3/eth"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TransferInfo holds the basic Transfer event data
type TransferInfo struct {
	From    common.Address
	To      common.Address
	TokenID string
}

const erc165ABI = `[{"constant":true,"inputs":[{"name":"interfaceId","type":"bytes4"}],"name":"supportsInterface","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"}]`

// getNFTHoldings retrieves current NFT holdings of a target address
func getNFTHoldings(client *ethclient.Client, targetAddress common.Address, startBlock *big.Int) (map[string]map[string]bool, error) {
	eventSignature := []byte("Transfer(address,address,uint256)")
	eventID := crypto.Keccak256Hash(eventSignature)

	type result struct {
		logs []types.Log
		err  error
	}

	logsChan := make(chan result, 2)

	// Query where 'from' is targetAddress
	go func() {
		query := ethereum.FilterQuery{
			FromBlock: startBlock,
			Topics: [][]common.Hash{
				{eventID},
				{common.BytesToHash(targetAddress.Bytes())},
			},
		}
		logs, err := client.FilterLogs(context.Background(), query)
		logsChan <- result{logs: logs, err: err}
	}()

	// Query where 'to' is targetAddress
	go func() {
		query := ethereum.FilterQuery{
			FromBlock: startBlock,
			Topics: [][]common.Hash{
				{eventID},
				nil,
				{common.BytesToHash(targetAddress.Bytes())},
			},
		}
		logs, err := client.FilterLogs(context.Background(), query)
		logsChan <- result{logs: logs, err: err}
	}()

	var allLogs []types.Log
	for range 2 {
		res := <-logsChan
		if res.err != nil {
			return nil, res.err
		}
		allLogs = append(allLogs, res.logs...)
	}

	// Sort logs by block number to ensure correct order
	sort.Slice(allLogs, func(i, j int) bool {
		return allLogs[i].BlockNumber < allLogs[j].BlockNumber
	})

	nftHoldings := make(map[string]map[string]bool)

	for _, vLog := range allLogs {
		from := common.HexToAddress(vLog.Topics[1].Hex())
		to := common.HexToAddress(vLog.Topics[2].Hex())
		tokenID := new(big.Int).SetBytes(vLog.Topics[3].Bytes()).String()
		contract := vLog.Address.Hex()

		if _, ok := nftHoldings[contract]; !ok {
			nftHoldings[contract] = make(map[string]bool)
		}

		if to == targetAddress {
			nftHoldings[contract][tokenID] = true
		}
		if from == targetAddress {
			delete(nftHoldings[contract], tokenID)
		}
	}

	return nftHoldings, nil
}

func isERC721(client *ethclient.Client, contractAddr string) (bool, error) {
	// Convert the contract address from string to Address
	contractAddress := common.HexToAddress(contractAddr)

	// Parse the ERC-165 ABI
	contractABI, err := abi.JSON(strings.NewReader(erc165ABI))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	// ERC-721 interfaceId is "0x80ac58cd"
	interfaceIdBytes := [4]byte{0x80, 0xac, 0x58, 0xcd}

	// Debug: Print the value of interfaceIdBytes
	fmt.Printf("InterfaceIdBytes: %x\n", interfaceIdBytes)

	// Pack the "supportsInterface" method with the interfaceIdBytes as argument
	callData, err := contractABI.Pack("supportsInterface", interfaceIdBytes)
	if err != nil {
		log.Fatalf("Failed to pack data: %v", err)
	}

	// Debug: Print the packed call data
	fmt.Printf("CallData: %x\n", callData)

	// Send the request to the contract and get the result
	result, err := client.CallContract(context.Background(), ethereum.CallMsg{
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

// Print out all NFTs currently held by the target address
func printHoldings(client *ethclient.Client, targetAddress common.Address, holdings map[string]map[string]bool) {
	fmt.Printf("NFTs currently held by %s:\n", targetAddress.Hex())
	for contract, tokens := range holdings {
		if len(tokens) == 0 {
			continue
		}
		fmt.Printf("- Contract: %s\n", contract)
		is, err := isERC721(client, contract)
		if err != nil {
			fmt.Println("Failed to check if contract is ERC721:", err)
		}
		fmt.Println("is ERC721: ", is)
		for tokenID := range tokens {
			fmt.Printf("    - TokenID: %s\n", tokenID)
		}
	}
}

func main() {
	// client, err := ethclient.Dial("wss://sepolia.infura.io/ws/v3/da05d3dc31244bd483a28d746233d32f")
	// if err != nil {
	// 	log.Fatalf("Failed to connect to Ethereum node: %v", err)
	// }

	// targetAddress := common.HexToAddress("0x7556989c2A60E60F0c66A2b9D77079BC9F189037")
	// startBlock := big.NewInt(0) // Adjust the starting block if needed

	// start := time.Now()
	// holdings, err := getNFTHoldings(client, targetAddress, startBlock)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// elapsed := time.Since(start)
	// fmt.Printf("Time taken: %s\n", elapsed)

	// // Display the result
	// printHoldings(client, targetAddress, holdings)

	eth := eth.NewEth("wss://sepolia.infura.io/ws/v3/da05d3dc31244bd483a28d746233d32f")

	defer eth.Provider.Close()

	start := time.Now()
	collections, err := eth.NFT.GetWalletNFTs("0xa84e540D1eb5458DFC2bC25760bD64fbECb8e345")

	if err != nil {
		panic(err)
	}

	elapsed := time.Since(start)
	fmt.Printf("Time taken: %s\n", elapsed)

	for _, collection := range collections {
		fmt.Println("collection address:", collection.ContractAddress)

		for _, nft := range collection.Tokens {
			fmt.Println("token id:", nft.TokenId)
		}
	}

}
