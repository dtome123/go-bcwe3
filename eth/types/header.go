package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Header struct {
	Origin           *types.Header `json:"-"`
	ParentHash       string        `json:"parent_hash"`
	UncleHash        string        `json:"sha3_uncles"`
	Coinbase         string        `json:"miner"`
	Root             string        `json:"state_root"`
	TxHash           string        `json:"transactions_root"`
	ReceiptHash      string        `json:"receipts_root"`
	Difficulty       *big.Int      `json:"difficulty"`
	Number           *big.Int      `json:"number"`
	GasLimit         uint64        `json:"gas_limit"`
	GasUsed          uint64        `json:"gas_used"`
	Time             uint64        `json:"timestamp"`
	Extra            []byte        `json:"extra_data"`
	MixDigest        string        `json:"mix_hash"`
	Nonce            uint64        `json:"nonce"`
	BaseFee          *big.Int      `json:"base_fee_per_gas"`
	WithdrawalsHash  *common.Hash  `json:"withdrawals_root"`
	BlobGasUsed      *uint64       `json:"blob_gas_used"`
	ExcessBlobGas    *uint64       `json:"excess_blob_gas"`
	ParentBeaconRoot *common.Hash  `json:"parent_beacon_block_root"`
	RequestsHash     *common.Hash  `json:"requests_hash"`
}

func WrapHeader(header *types.Header) *Header {

	if header == nil {
		return nil
	}

	return &Header{
		Origin:      header,
		ParentHash:  header.ParentHash.String(),
		UncleHash:   header.UncleHash.String(),
		Coinbase:    header.Coinbase.String(),
		Root:        header.Root.String(),
		TxHash:      header.TxHash.String(),
		ReceiptHash: header.ReceiptHash.String(),
		Difficulty:  header.Difficulty,
		Number:      header.Number,
		GasLimit:    header.GasLimit,
		GasUsed:     header.GasUsed,
		Time:        header.Time,
		Extra:       header.Extra,
		MixDigest:   header.MixDigest.String(),
		Nonce:       header.Nonce.Uint64(),
		BaseFee:     header.BaseFee,

		// retaining the exact types of EIP-related fields to maintain compliance with canonical chain data formats.
		WithdrawalsHash:  header.WithdrawalsHash,
		BlobGasUsed:      header.BlobGasUsed,
		ExcessBlobGas:    header.ExcessBlobGas,
		ParentBeaconRoot: header.ParentBeaconRoot,
		RequestsHash:     header.RequestsHash,
	}
}
