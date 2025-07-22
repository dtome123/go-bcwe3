package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"

	"github.com/dtome123/go-bcwe3/eth"
	"github.com/ethereum/go-ethereum/core/types"
)

func main() {
	infuraURL := "https://sepolia.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161"
	var signedTx *types.Transaction
	var buf bytes.Buffer
	signedTx.EncodeRLP(&buf)
	rawTxHex := hex.EncodeToString(buf.Bytes())

	eth := eth.NewEth(infuraURL)

	finalTx, err := eth.GetProvider().SendSignedTransaction(context.Background(), rawTxHex)
	if err != nil {
		panic(err)
	}

	fmt.Println("tx hash: ", finalTx)
}
