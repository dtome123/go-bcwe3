package provider

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/dtome123/go-bcwe3/eth/types"
	"github.com/dtome123/go-bcwe3/eth/utils"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	goethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
)

type impl struct {
	client *ethclient.Client
}

func NewProvider(dsn string) Provider {
	client, err := ethclient.Dial(dsn)
	if err != nil {
		log.Fatal(err)
	}

	c := &impl{
		client: client,
	}

	return c
}

func (e *impl) Close() {
	e.client.Close()
}

////////////////////// wrapper methods for ethclient //////////////////////

func (e *impl) Client() *rpc.Client {
	return e.client.Client()
}

func (e *impl) ChainID(ctx context.Context) (*big.Int, error) {
	return e.client.ChainID(ctx)
}

func (e *impl) BlockByHash(ctx context.Context, hash string) (*types.Block, error) {
	block, err := e.client.BlockByHash(ctx, common.HexToHash(hash))
	if err != nil {
		return nil, err
	}

	return types.WrapBlock(block), nil
}

func (e *impl) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	block, err := e.client.BlockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	return types.WrapBlock(block), nil
}

func (e *impl) BlockNumber(ctx context.Context) (uint64, error) {
	blockNumber, err := e.client.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}
	return blockNumber, nil
}

func (e *impl) PeerCount(ctx context.Context) (uint64, error) {
	return e.client.PeerCount(ctx)
}

func (e *impl) BlockReceipts(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) ([]*types.Receipt, error) {

	receipts, err := e.client.BlockReceipts(ctx, blockNrOrHash)
	if err != nil {
		return nil, err
	}

	return types.WrapReceipts(receipts), nil
}

func (e *impl) HeaderByHash(ctx context.Context, hash string) (*types.Header, error) {
	header, err := e.client.HeaderByHash(ctx, common.HexToHash(hash))

	return types.WrapHeader(header), err
}

func (e *impl) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	header, err := e.client.HeaderByNumber(ctx, number)

	return types.WrapHeader(header), err
}

func (e *impl) TransactionByHash(ctx context.Context, hash string) (tx *types.CompleteTx, isPending bool, err error) {

	transaction, isPending, err := e.client.TransactionByHash(ctx, common.HexToHash(hash))
	if err != nil {
		return
	}

	tx, err = e.buildCompleteTransactionWithBlockHash(transaction.Hash(), transaction)

	return
}

func (e *impl) TransactionSender(ctx context.Context, tx *types.CompleteTx, blockHash string, index uint) (common.Address, error) {
	return e.client.TransactionSender(ctx, tx.Origin, common.HexToHash(blockHash), index)
}

func (e *impl) TransactionCount(ctx context.Context, blockHash string) (uint, error) {
	return e.client.TransactionCount(ctx, common.HexToHash(blockHash))
}

func (e *impl) TransactionInBlock(ctx context.Context, blockHash string, index uint) (*types.CompleteTx, error) {

	tx, err := e.client.TransactionInBlock(ctx, common.HexToHash(blockHash), index)
	if err != nil {
		return nil, err
	}

	return e.buildCompleteTransactionWithBlockHash(common.HexToHash(blockHash), tx)
}

func (e *impl) TransactionReceipt(ctx context.Context, hash string) (*types.Receipt, error) {
	receipt, err := e.client.TransactionReceipt(ctx, common.HexToHash(hash))

	return types.WrapReceipt(receipt), err
}

func (e *impl) BalanceAt(ctx context.Context, account string, blockNumber *big.Int) (*big.Int, error) {

	address := common.HexToAddress(account)

	return e.client.BalanceAt(ctx, address, blockNumber)
}

func (e *impl) BalanceAtHash(ctx context.Context, account string, blockHash string) (*big.Int, error) {
	accountAddress := common.HexToAddress(account)
	return e.client.BalanceAtHash(ctx, accountAddress, common.HexToHash(blockHash))
}

func (e *impl) StorageAt(ctx context.Context, account string, key string, blockNumber *big.Int) ([]byte, error) {
	accountAddress := common.HexToAddress(account)
	return e.client.StorageAt(ctx, accountAddress, common.HexToHash(key), blockNumber)
}

func (e *impl) StorageAtHash(ctx context.Context, account string, key string, blockHash string) ([]byte, error) {
	accountAddress := common.HexToAddress(account)
	return e.client.StorageAtHash(ctx, accountAddress, common.HexToHash(key), common.HexToHash(blockHash))
}

func (e *impl) CodeAt(ctx context.Context, account string, blockNumber *big.Int) ([]byte, error) {
	accountAddress := common.HexToAddress(account)
	return e.client.CodeAt(ctx, accountAddress, blockNumber)
}

func (e *impl) CodeAtHash(ctx context.Context, account string, blockHash string) ([]byte, error) {
	accountAddress := common.HexToAddress(account)
	return e.client.CodeAtHash(ctx, accountAddress, common.HexToHash(blockHash))
}

func (e *impl) NonceAt(ctx context.Context, account string, blockNumber *big.Int) (uint64, error) {
	accountAddress := common.HexToAddress(account)
	return e.client.NonceAt(ctx, accountAddress, blockNumber)
}

func (e *impl) NonceAtHash(ctx context.Context, account string, blockHash string) (uint64, error) {
	accountAddress := common.HexToAddress(account)
	return e.client.NonceAtHash(ctx, accountAddress, common.HexToHash(blockHash))
}
func (e *impl) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	logs, err := e.client.FilterLogs(ctx, q)

	wLogs := make([]types.Log, len(logs))
	for i, l := range logs {
		p := types.WrapLog(&l)
		wLogs[i] = p.Dereference()
	}

	return wLogs, err
}
func (e *impl) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {

	oCh := make(chan goethTypes.Log)

	sub, err := e.client.SubscribeFilterLogs(ctx, q, oCh)
	if err != nil {
		close(oCh)
		return nil, err
	}

	go func() {
		defer close(oCh)
		for l := range oCh {
			wLog := types.WrapLog(&l)
			ch <- wLog.Dereference()
		}
	}()

	return sub, nil
}

func (e *impl) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	oCh := make(chan *goethTypes.Header)

	sub, err := e.client.SubscribeNewHead(ctx, oCh)
	if err != nil {
		close(oCh)
		return nil, err
	}

	go func() {
		defer close(oCh)
		for h := range oCh {
			wHeader := types.WrapHeader(h)
			ch <- wHeader
		}
	}()

	return sub, nil

}

func (e *impl) PendingBalanceAt(ctx context.Context, account string) (*big.Int, error) {
	accountAddress := common.HexToAddress(account)
	return e.client.PendingBalanceAt(ctx, accountAddress)
}
func (e *impl) PendingStorageAt(ctx context.Context, account string, key string) ([]byte, error) {
	accountAddress := common.HexToAddress(account)
	return e.client.PendingStorageAt(ctx, accountAddress, common.HexToHash(key))
}

func (e *impl) PendingCodeAt(ctx context.Context, account string) ([]byte, error) {
	accountAddress := common.HexToAddress(account)
	return e.client.PendingCodeAt(ctx, accountAddress)
}
func (e *impl) PendingNonceAt(ctx context.Context, account string) (uint64, error) {
	accountAddress := common.HexToAddress(account)
	return e.client.PendingNonceAt(ctx, accountAddress)
}
func (e *impl) PendingTransactionCount(ctx context.Context) (uint, error) {
	return e.client.PendingTransactionCount(ctx)
}
func (e *impl) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return e.client.CallContract(ctx, msg, blockNumber)
}
func (e *impl) CallContractAtHash(ctx context.Context, msg ethereum.CallMsg, blockHash string) ([]byte, error) {
	return e.client.CallContractAtHash(ctx, msg, common.HexToHash(blockHash))
}
func (e *impl) PendingCallContract(ctx context.Context, msg ethereum.CallMsg) ([]byte, error) {
	return e.client.PendingCallContract(ctx, msg)
}
func (e *impl) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return e.client.SuggestGasPrice(ctx)
}
func (e *impl) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return e.client.SuggestGasTipCap(ctx)
}
func (e *impl) BlobBaseFee(ctx context.Context) (*big.Int, error) {
	return e.client.BlobBaseFee(ctx)
}
func (e *impl) FeeHistory(ctx context.Context, blockCount uint64, lastBlock *big.Int, rewardPercentiles []float64) (*ethereum.FeeHistory, error) {
	return e.client.FeeHistory(ctx, blockCount, lastBlock, rewardPercentiles)
}
func (e *impl) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	return e.client.EstimateGas(ctx, msg)
}
func (e *impl) EstimateGasAtBlock(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) (uint64, error) {
	return e.client.EstimateGasAtBlock(ctx, msg, blockNumber)
}
func (e *impl) EstimateGasAtBlockHash(ctx context.Context, msg ethereum.CallMsg, blockHash string) (uint64, error) {
	return e.client.EstimateGasAtBlockHash(ctx, msg, common.HexToHash(blockHash))
}
func (e *impl) SendTransaction(ctx context.Context, tx *goethTypes.Transaction) error {
	return e.client.SendTransaction(ctx, tx)
}

//////////////////////////////// EXTRA ////////////////////////////////

func (e *impl) CalculateTxFee(tx *types.Tx) (*big.Int, error) {

	receipt, err := e.client.TransactionReceipt(context.Background(), tx.Origin.Hash())
	if err != nil {
		return nil, err
	}

	if tx == nil || receipt == nil {
		return big.NewInt(0), fmt.Errorf("tx or receipt is nil")
	}

	gasUsed := new(big.Int).SetUint64(receipt.GasUsed)
	gasPrice := tx.Origin.GasPrice()

	fee := new(big.Int).Mul(gasUsed, gasPrice)
	return fee, nil
}

func (e *impl) SendSignedTransaction(ctx context.Context, signedTxHex string) (string, error) {
	data := common.FromHex(signedTxHex)

	var tx goethTypes.Transaction
	if err := rlp.DecodeBytes(data, &tx); err != nil {
		log.Fatal(err)
	}

	err := e.client.SendTransaction(ctx, &tx)

	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

func (e *impl) IsBlockFinalized(ctx context.Context, blockNumber *big.Int) (bool, error) {
	// Get the finalized block using raw RPC call
	var finalizedBlock *types.Header
	err := e.client.Client().CallContext(ctx, &finalizedBlock, "eth_getBlockByNumber", "finalized", false)
	if err != nil {
		return false, err
	}

	return blockNumber.Cmp(finalizedBlock.Number) <= 0, nil
}

func (e *impl) BuildCompleteTransaction(block *types.Block, tx *types.Tx) (*types.CompleteTx, error) {
	return e.buildCompleteTransaction(block.Origin, tx.Origin)
}

// ////////////////////////// private ////////////////////////////////

func (e *impl) buildCompleteTransactionWithBlockHash(blockHash common.Hash, tx *goethTypes.Transaction) (*types.CompleteTx, error) {
	ctx := context.Background()

	block, _ := e.client.BlockByHash(ctx, blockHash)

	return e.buildCompleteTransaction(block, tx)
}

func (e *impl) buildCompleteTransaction(block *goethTypes.Block, tx *goethTypes.Transaction) (*types.CompleteTx, error) {

	ctx := context.Background()

	receipt, _ := e.client.TransactionReceipt(ctx, tx.Hash())
	wrapTx := types.WrapTx(tx)

	var timestamp uint64
	if block != nil {
		timestamp = block.Time()
	} else {
		blockObj, err := e.client.BlockByHash(context.Background(), receipt.BlockHash)
		if err == nil {
			timestamp = blockObj.Time()
		}
	}

	from := utils.GetFromAddressTx(tx)
	to := utils.GetToAddressTx(tx)
	fee, _ := e.CalculateTxFee(wrapTx)

	complete := &types.CompleteTx{
		Origin:    tx,
		Hash:      tx.Hash().Hex(),
		From:      from,
		To:        to,
		Value:     tx.Value(),
		Gas:       tx.Gas(),
		GasPrice:  tx.GasPrice(),
		GasUsed:   receipt.GasUsed,
		Fee:       fee,
		Nonce:     tx.Nonce(),
		Status:    receipt.Status,
		BlockHash: receipt.BlockHash.Hex(),
		BlockNum:  receipt.BlockNumber.Uint64(),
		Timestamp: timestamp,
		Pending:   receipt.Status == 0 && receipt.BlockNumber == nil,
		GasFee: &types.GasFee{
			BaseFee: block.BaseFee(),
			TipCap:  tx.GasTipCap(),
			FeeCap:  tx.GasFeeCap(),
		},
		Type: tx.Type(),
	}

	if tx.To() != nil {
		complete.To = tx.To().Hex()
	}

	return complete, nil
}
