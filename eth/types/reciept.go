package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
)

type Receipt struct {
	Origin            *types.Receipt `json:"-"`
	Type              uint8          `json:"type"`
	PostState         []byte         `json:"root"`
	Status            uint64         `json:"status"`
	CumulativeGasUsed uint64         `json:"cumulative_gas_used"`
	Logs              []*Log         `json:"logs"`
	TxHash            string         `json:"tx_hash"`
	ContractAddress   string         `json:"contract_address"`
	GasUsed           uint64         `json:"gas_used"`
	EffectiveGasPrice *big.Int       `json:"effective_gas_price"`
	BlobGasUsed       uint64         `json:"blob_gas_used,omitempty"`
	BlobGasPrice      *big.Int       `json:"blob_gas_price,omitempty"`
	BlockHash         string         `json:"block_hash,omitempty"`
	BlockNumber       *big.Int       `json:"block_number,omitempty"`
	TransactionIndex  uint           `json:"transaction_index"`
}

func WrapReceipt(r *types.Receipt) *Receipt {

	if r == nil {
		return nil
	}

	return &Receipt{
		Origin:            r,
		Type:              r.Type,
		PostState:         r.PostState,
		Status:            r.Status,
		CumulativeGasUsed: r.CumulativeGasUsed,
		Logs:              WrapLogs(r.Logs),
		TxHash:            r.TxHash.String(),
		ContractAddress:   r.ContractAddress.String(),
		GasUsed:           r.GasUsed,
		EffectiveGasPrice: r.EffectiveGasPrice,
		BlobGasUsed:       r.BlobGasUsed,
		BlobGasPrice:      r.BlobGasPrice,
		BlockHash:         r.BlockHash.String(),
		BlockNumber:       r.BlockNumber,
		TransactionIndex:  r.TransactionIndex,
	}
}

func WrapReceipts(receipts []*types.Receipt) []*Receipt {
	wrappedReceipts := make([]*Receipt, len(receipts))
	for i, receipt := range receipts {
		wrappedReceipts[i] = WrapReceipt(receipt)
	}
	return wrappedReceipts
}