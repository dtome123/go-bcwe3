package provider

import (
	"context"
	"math/big"

	"github.com/dtome123/go-bcwe3/eth/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	goethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

type Provider interface {
	Close()
	Client() *rpc.Client
	ChainID(ctx context.Context) (*big.Int, error)
	BlockByHash(ctx context.Context, hash string) (*types.Block, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	BlockNumber(ctx context.Context) (uint64, error)
	PeerCount(ctx context.Context) (uint64, error)
	BlockReceipts(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) ([]*goethTypes.Receipt, error)
	HeaderByHash(ctx context.Context, hash string) (*goethTypes.Header, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*goethTypes.Header, error)
	TransactionByHash(ctx context.Context, hash string) (tx *types.CompleteTx, isPending bool, err error)
	TransactionSender(ctx context.Context, tx *types.CompleteTx, block string, index uint) (common.Address, error)
	TransactionCount(ctx context.Context, blockHash string) (uint, error)
	TransactionInBlock(ctx context.Context, blockHash string, index uint) (*types.CompleteTx, error)
	TransactionReceipt(ctx context.Context, txHash string) (*goethTypes.Receipt, error)
	BalanceAt(ctx context.Context, account string, blockNumber *big.Int) (*big.Int, error)
	BalanceAtHash(ctx context.Context, account string, blockHash string) (*big.Int, error)
	StorageAt(ctx context.Context, account string, key string, blockNumber *big.Int) ([]byte, error)
	StorageAtHash(ctx context.Context, account string, key string, blockHash string) ([]byte, error)
	CodeAt(ctx context.Context, account string, blockNumber *big.Int) ([]byte, error)
	CodeAtHash(ctx context.Context, account string, blockHash string) ([]byte, error)
	NonceAt(ctx context.Context, account string, blockNumber *big.Int) (uint64, error)
	NonceAtHash(ctx context.Context, account string, blockHash string) (uint64, error)
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]goethTypes.Log, error)
	SubscribeNewHead(ctx context.Context, ch chan<- *goethTypes.Header) (ethereum.Subscription, error)
	SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- goethTypes.Log) (ethereum.Subscription, error)
	PendingBalanceAt(ctx context.Context, account string) (*big.Int, error)
	PendingStorageAt(ctx context.Context, account string, key string) ([]byte, error)
	PendingCodeAt(ctx context.Context, account string) ([]byte, error)
	PendingNonceAt(ctx context.Context, account string) (uint64, error)
	PendingTransactionCount(ctx context.Context) (uint, error)
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	CallContractAtHash(ctx context.Context, msg ethereum.CallMsg, blockHash string) ([]byte, error)
	PendingCallContract(ctx context.Context, msg ethereum.CallMsg) ([]byte, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	SuggestGasTipCap(ctx context.Context) (*big.Int, error)
	BlobBaseFee(ctx context.Context) (*big.Int, error)
	FeeHistory(ctx context.Context, blockCount uint64, lastBlock *big.Int, rewardPercentiles []float64) (*ethereum.FeeHistory, error)
	EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
	EstimateGasAtBlock(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) (uint64, error)
	EstimateGasAtBlockHash(ctx context.Context, msg ethereum.CallMsg, blockHash string) (uint64, error)
	SendTransaction(ctx context.Context, tx *goethTypes.Transaction) error

	// extra
	CalculateTxFee(tx *types.Tx) (*big.Int, error)
	SendSignedTransaction(ctx context.Context, signedTxHex string) error
	IsBlockFinalized(ctx context.Context, blockNumber *big.Int) (bool, error)
}
