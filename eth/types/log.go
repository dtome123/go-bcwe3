package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Log struct {
	Origin      *types.Log    `json:"-"`
	Address     string        `json:"address"`
	Topics      []common.Hash `json:"topics"`
	Data        []byte        `json:"data"`
	BlockNumber uint64        `json:"block_number"`
	TxHash      string        `json:"transaction_hash"`
	TxIndex     uint          `json:"transaction_index"`
	BlockHash   string        `json:"block_hash"`
	Index       uint          `json:"log_index"`
	Removed     bool          `json:"removed"`
}

func WrapLog(log *types.Log) *Log {

	if log == nil {
		return nil
	}

	return &Log{
		Origin:      log,
		Address:     log.Address.String(),
		Topics:      log.Topics,
		Data:        log.Data,
		BlockNumber: log.BlockNumber,
		TxHash:      log.TxHash.String(),
		TxIndex:     log.TxIndex,
		BlockHash:   log.BlockHash.String(),
		Index:       log.Index,
		Removed:     log.Removed,
	}
}

func WrapLogs(logs []*types.Log) []*Log {
	wrappedLogs := make([]*Log, len(logs))
	for i, log := range logs {
		wrappedLogs[i] = WrapLog(log)
	}
	return wrappedLogs
}

func (l *Log) Dereference() Log {
	return *l
}
