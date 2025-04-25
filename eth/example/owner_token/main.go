package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Định nghĩa sự kiện Transfer của ERC721
type TransferEvent struct {
	From    common.Address
	To      common.Address
	TokenID *big.Int
}

func main() {
	// URL của Ethereum node (Infura hoặc địa chỉ node của bạn)
	infuraURL := "wss://sepolia.infura.io/ws/v3/da05d3dc31244bd483a28d746233d32f"
	contractAddress := "0x975Bda1d9287433B868967685bC79637A08EfBEc"
	// Kết nối với Ethereum node
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Không thể kết nối với node: %v", err)
	}

	// Chuyển đổi contract address thành common.Address
	address := common.HexToAddress(contractAddress)

	// Lắng nghe sự kiện Transfer (ERC721 có sự kiện Transfer)
	eventSignature := []byte("Transfer(address,address,uint256)")
	eventID := crypto.Keccak256Hash(eventSignature)

	// Tạo filter log để lắng nghe sự kiện Transfer
	query := ethereum.FilterQuery{
		Addresses: []common.Address{address},
		Topics:    [][]common.Hash{{eventID}},
	}

	// Đăng ký lắng nghe sự kiện Transfer
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatalf("Lỗi khi đăng ký lắng nghe sự kiện: %v", err)
	}

	sort.Slice(logs, func(i, j int) bool {
		return logs[i].BlockNumber < logs[j].BlockNumber
	})

	// Lặp qua các logs và giải mã sự kiện Transfer
	owners := make(map[string]string)

	for idx, log := range logs {

		fmt.Println("log ", idx)
		for i := range log.Topics {
			fmt.Println("log.Topics[i].Hex", log.Topics[i].Hex())
		}

		if len(log.Topics) != 4 {
			continue
		}

		from := common.HexToAddress(log.Topics[1].Hex())
		to := common.HexToAddress(log.Topics[2].Hex())
		tokenId := new(big.Int).SetBytes(log.Topics[3].Bytes()).String()

		if to == (common.Address{}) {
			// Bị burn
			delete(owners, tokenId)
			fmt.Printf("Token %s was burned (from %s)\n", tokenId, from.Hex())
		} else {
			owners[tokenId] = to.Hex()
			fmt.Printf("Token %s now owned by %s (from %s)\n", tokenId, to.Hex(), from.Hex())
		}
	}

	// In toàn bộ owner hiện tại
	fmt.Println("\n== Owner hiện tại của các Token ==")
	for tokenId, owner := range owners {
		fmt.Printf("TokenID %s → Owner: %s\n", tokenId, owner)
	}
}
