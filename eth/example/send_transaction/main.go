package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/dtome123/go-bcwe3/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// URL của Ethereum node (Infura hoặc địa chỉ node của bạn)
	infuraURL := "wss://sepolia.infura.io/ws/v3/da05d3dc31244bd483a28d746233d32f"

	// Kết nối với Ethereum node
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Không thể kết nối với node: %v", err)
	}

	privateKey, err := crypto.HexToECDSA("<private key>")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(1000000000000) // 1
	gasLimit := uint64(21000)          // gas limit for standard tx
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	toAddress := common.HexToAddress("0xa84e540D1eb5458DFC2bC25760bD64fbECb8e345")
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(chainID)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// err = client.SendTransaction(context.Background(), signedTx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println("Signed tx hash:", signedTx.Hash().Hex())

	var buf bytes.Buffer
	signedTx.EncodeRLP(&buf)
	rawTxHex := hex.EncodeToString(buf.Bytes())
	fmt.Println("Raw signed tx:", "0x"+rawTxHex)

	eth := eth.NewEth(infuraURL)

	finalTx, err := eth.Provider().SendSignedTransaction(context.Background(), rawTxHex)

	if err != nil {
		panic(err)
	}

	fmt.Println("tx hash: ", finalTx)
}
