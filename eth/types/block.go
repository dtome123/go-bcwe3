package types

import (
	"math/big"

	"github.com/dtome123/go-bcwe3/eth/utils"
	"github.com/ethereum/go-ethereum/core/types"
)

type Tx struct {
	Origin   *types.Transaction `json:"-"`
	Hash     string             `json:"hash"`
	From     string             `json:"from"`
	To       string             `json:"to"`
	Value    *big.Int           `json:"value"`
	Gas      uint64             `json:"gas"`
	GasPrice *big.Int           `json:"gas_price"`
}

type CompleteTx struct {
	Origin    *types.Transaction `json:"-"`
	Hash      string             `json:"hash"`
	From      string             `json:"from"`
	To        string             `json:"to,omitempty"`
	Value     *big.Int           `json:"value"`
	Gas       uint64             `json:"gas"`
	GasPrice  *big.Int           `json:"gas_price"`
	GasUsed   uint64             `json:"gas_used"`
	GasFee    *GasFee            `json:"gas_fee"`
	Fee       *big.Int           `json:"fee"`
	Nonce     uint64             `json:"nonce"`
	Status    uint64             `json:"status"`
	BlockHash string             `json:"block_hash"`
	BlockNum  uint64             `json:"block_number"`
	Timestamp uint64             `json:"timestamp"`
	Pending   bool               `json:"pending"`
	Type      uint8              `json:"type"`
}

type GasFee struct {
	BaseFee *big.Int `json:"base_fee"`
	TipCap  *big.Int `json:"tip_cap"`
	FeeCap  *big.Int `json:"fee_cap"`
}

type Block struct {
	Origin       *types.Block  `json:"-"`
	Number       *big.Int      `json:"block_number"`
	Hash         string        `json:"block_hash"`
	ParentHash   string        `json:"parent_hash"`
	Nonce        uint64        `json:"nonce"`
	Time         uint64        `json:"timestamp"`
	Miner        string        `json:"miner"`
	GasLimit     uint64        `json:"gas_limit"`
	GasUsed      uint64        `json:"gas_used"`
	Transactions []*Tx         `json:"transactions"`
	Withdrawals  []*Withdrawal `json:"withdrawals"`
}

type Withdrawal struct {
	Index     uint64 `json:"index"`
	Validator uint64 `json:"validatorIndex"`
	Address   string `json:"address"`
	Amount    uint64 `json:"amount"`
}

func WrapBlock(block *types.Block) *Block {
	var txs []*Tx
	for _, tx := range block.Transactions() {
		txs = append(txs, WrapTx(tx))
	}

	var withdrawals []*Withdrawal
	for _, withdrawal := range block.Withdrawals() {
		withdrawals = append(withdrawals, &Withdrawal{
			Index:     withdrawal.Index,
			Validator: withdrawal.Validator,
			Address:   withdrawal.Address.Hex(),
			Amount:    withdrawal.Amount,
		})
	}

	blockInfo := &Block{
		Origin:       block,
		Number:       block.Number(),
		Hash:         block.Hash().Hex(),
		ParentHash:   block.ParentHash().Hex(),
		Nonce:        block.Nonce(),
		Time:         block.Time(),
		Miner:        block.Coinbase().Hex(),
		GasLimit:     block.GasLimit(),
		GasUsed:      block.GasUsed(),
		Transactions: txs,
		Withdrawals:  withdrawals,
	}

	return blockInfo
}

func WrapTx(tx *types.Transaction) *Tx {
	from := utils.GetFromAddressTx(tx)
	to := utils.GetToAddressTx(tx)

	return &Tx{
		Origin:   tx,
		Hash:     tx.Hash().Hex(),
		From:     from,
		To:       to,
		Value:    tx.Value(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
	}
}
