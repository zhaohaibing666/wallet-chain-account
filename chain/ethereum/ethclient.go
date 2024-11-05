package ethereum

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/url"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/zhaohaibing666/wallet-chain-account/common/global_const"
	"github.com/zhaohaibing666/wallet-chain-account/common/retry"
)

const (
	defaultDialTimeout    = 5 * time.Second
	defaultDialAttempts   = 5
	defaultRequestTimeout = 10 * time.Second
)

type TransactionList struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Hash  string `json:"hash"`
	Value string `json:"value"`
}

type RpcBlock struct {
	Hash         common.Hash       `json:"hash"`
	Height       uint64            `json:"height"`
	Transactions []TransactionList `json:"transactions"`
	BaseFee      string            `json:"baseFeePerGas"`
}

type EthClient interface {
	BlockHeaderByNumber(*big.Int) (*types.Header, error)

	BlockByNumber(*big.Int) (*RpcBlock, error)
	BlockByHash(common.Hash) (*RpcBlock, error)

	LatestSafeBlockHeader() (*types.Header, error)
	LatestFinalizedBlockHeader() (*types.Header, error)
	BlockHeaderByHash(common.Hash) (*types.Header, error)
	BlockHeadersByRange(*big.Int, *big.Int, uint) ([]types.Header, error)

	TxByHash(common.Hash) (*types.Transaction, error)
	TxReceiptByHash(common.Hash) (*types.Receipt, error)

	StorageHash(common.Address, *big.Int) (common.Hash, error)
	FilterLogs(filterQuery ethereum.FilterQuery, chainId uint) (Logs, error)

	TxCountByAddress(common.Address) (hexutil.Uint64, error)

	SendRawTransaction(rawTx string) error

	SuggestGasPrice() (*big.Int, error)
	SuggestGasTipCap() (*big.Int, error)

	Close()
}

type Clnt struct {
	rpc RPC
}

// BlockByHash implements EthClient.
func (c *Clnt) BlockByHash(hash common.Hash) (*RpcBlock, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()
	var block *RpcBlock
	err := c.rpc.CallContext(ctxwt, &block, "eth_getBlockByHash", hash, true)
	if err != nil {
		return nil, err
	}
	if block == nil {
		return nil, ethereum.NotFound
	}

	return block, nil
}

// BlockByNumber implements EthClient.
func (c *Clnt) BlockByNumber(number *big.Int) (*RpcBlock, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()
	var block *RpcBlock
	err := c.rpc.CallContext(ctxwt, &block, "eth_getBlockByNumber", toBlockNumArg(number), true)
	if err != nil {
		return nil, err
	}
	if block == nil {
		return nil, ethereum.NotFound
	}

	return block, nil
}

// BlockHeaderByHash implements EthClient.
func (c *Clnt) BlockHeaderByHash(hash common.Hash) (*types.Header, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	var header *types.Header
	err := c.rpc.CallContext(ctxwt, &header, "eth_getBlockByHash", hash, false)
	if err != nil {
		return nil, err
	}
	if header == nil {
		return nil, ethereum.NotFound
	}

	if header.Hash() != hash {
		return nil, errors.New("hash mismatch")
	}
	return header, nil
}

// BlockHeaderByNumber implements EthClient.
func (c *Clnt) BlockHeaderByNumber(number *big.Int) (*types.Header, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()
	var head *types.Header
	err := c.rpc.CallContext(ctxwt, &head, "eth_getBlockByNumber", toBlockNumArg(number), false)
	if err != nil {
		return nil, err
	}
	if head == nil {
		return nil, ethereum.NotFound
	}

	return head, nil
}

// BlockHeadersByRange implements EthClient.
func (c *Clnt) BlockHeadersByRange(startHeight *big.Int, endHeight *big.Int, chainId uint) ([]types.Header, error) {
	if startHeight.Cmp(endHeight) == 0 {
		header, err := c.BlockHeaderByNumber(startHeight)
		if err != nil {
			return nil, err
		}
		return []types.Header{*header}, nil
	}
	count := new(big.Int).Sub(endHeight, startHeight).Uint64() + 1
	headers := make([]types.Header, count)
	batchElems := make([]rpc.BatchElem, count)

	cxtwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	if chainId == uint(global_const.ZkFairSepoliaChainId) || chainId == uint(global_const.ZkFairChainId) {
		groupSize := 100
		var wg sync.WaitGroup
		numGroup := int(count-1)/groupSize + 1
		wg.Add(numGroup)

		for i := 0; i < int(count); i += groupSize {
			start := i
			end := i + groupSize - 1
			if end > int(count) {
				end = int(count) - 1
			}
			go func(start, end int) {
				defer wg.Done()
				for j := 0; j <= end; j++ {
					height := new(big.Int).Add(startHeight, new(big.Int).SetUint64(uint64(j)))
					batchElems[j] = rpc.BatchElem{
						Method: "eth_getBlockByNumber",
						Result: new(types.Header),
						Error:  nil,
					}
					header := new(types.Header)
					batchElems[j].Error = c.rpc.CallContext(cxtwt, header, batchElems[j].Method, toBlockNumArg(height), false)
					batchElems[j].Result = header
				}
			}(start, end)
		}
		wg.Wait()
	} else {
		for i := uint64(0); i < count; i++ {
			height := new(big.Int).Add(startHeight, new(big.Int).SetUint64(i))
			batchElems[i] = rpc.BatchElem{
				Method: "eth_getBlockByNumber",
				Args:   []interface{}{toBlockNumArg(height), false},
				Result: &headers[i],
			}
			err := c.rpc.BatchCallContext(cxtwt, batchElems)
			if err != nil {
				return nil, err
			}
		}
	}

	for index, batchElem := range batchElems {
		header, ok := batchElem.Result.(*types.Header)
		if !ok {
			continue
		}
		headers[index] = *header
	}

	return headers, nil
}

// Close implements EthClient.
func (c *Clnt) Close() {
	panic("unimplemented")
}

// FilterLogs implements EthClient.
func (c *Clnt) FilterLogs(query ethereum.FilterQuery, chainId uint) (Logs, error) {
	arg, err := toFilterArg(query)
	if err != nil {
		return Logs{}, err
	}

	var logs []types.Log
	var header types.Header

	batchElems := make([]rpc.BatchElem, 2)
	batchElems[0] = rpc.BatchElem{Method: "eth_getBlockByNumber", Args: []interface{}{toBlockNumArg(query.ToBlock), false}, Result: &header}
	batchElems[1] = rpc.BatchElem{Method: "eth_getLogs", Args: []interface{}{arg}, Result: &logs}
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout*10)
	defer cancel()
	if chainId == uint(global_const.ZkFairSepoliaChainId) ||
		chainId == uint(global_const.ZkFairChainId) {

		batchElems[0].Error = c.rpc.CallContext(ctxwt, &header, batchElems[0].Method, toBlockNumArg(query.ToBlock), false)
		batchElems[1].Error = c.rpc.CallContext(ctxwt, &logs, batchElems[1].Method, arg)
	} else {
		err = c.rpc.BatchCallContext(ctxwt, batchElems)
		if err != nil {
			return Logs{}, err
		}
	}

	if batchElems[0].Error != nil {
		return Logs{}, fmt.Errorf("unable to query for the `FilterQuery#ToBlock` header: %w", batchElems[0].Error)
	}
	if batchElems[1].Error != nil {
		return Logs{}, fmt.Errorf("unable to query logs: %w", batchElems[1].Error)
	}
	return Logs{Logs: logs, ToBlockHeader: &header}, nil
}

// LatestFinalizedBlockHeader implements EthClient.
func (c *Clnt) LatestFinalizedBlockHeader() (*types.Header, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()
	var header *types.Header
	err := c.rpc.CallContext(ctxwt, &header, "eth_getBlockByNumber", "finalized", false)
	if err != nil {
		return nil, err
	}
	if header == nil {
		return nil, ethereum.NotFound
	}
	return header, nil
}

// LatestSafeBlockHeader implements EthClient.
func (c *Clnt) LatestSafeBlockHeader() (*types.Header, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()
	var header *types.Header
	err := c.rpc.CallContext(ctxwt, &header, "eth_getBlockByNumber", "safe", false)
	if err != nil {
		return nil, err
	} else if header == nil {
		return nil, ethereum.NotFound
	}
	return header, nil
}

// SendRawTransaction implements EthClient.
func (c *Clnt) SendRawTransaction(rawTx string) error {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()
	if err := c.rpc.CallContext(ctxwt, nil, "eth_sendRawTransaction", rawTx); err != nil {
		return err
	}
	return nil
}

// StorageHash implements EthClient.
func (c *Clnt) StorageHash(address common.Address, blockNumber *big.Int) (common.Hash, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	proof := struct{ StorageHash common.Hash }{}
	err := c.rpc.CallContext(ctxwt, &proof, "eth_getProof", address, nil, toBlockNumArg(blockNumber))
	if err != nil {
		return common.Hash{}, err
	}

	return proof.StorageHash, nil
}

// SuggestGasPrice implements EthClient.
func (c *Clnt) SuggestGasPrice() (*big.Int, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()
	var hex hexutil.Big
	if err := c.rpc.CallContext(ctxwt, &hex, "eth_gasPrice"); err != nil {
		return nil, err
	}
	return (*big.Int)(&hex), nil
}

// SuggestGasTipCap implements EthClient.
func (c *Clnt) SuggestGasTipCap() (*big.Int, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()
	var hex hexutil.Big
	if err := c.rpc.CallContext(ctxwt, &hex, "eth_maxPriorityFeePerGas"); err != nil {
		return nil, err
	}
	return (*big.Int)(&hex), nil
}

// TxByHash implements EthClient.
func (c *Clnt) TxByHash(hash common.Hash) (*types.Transaction, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	var tx *types.Transaction
	err := c.rpc.CallContext(ctxwt, &tx, "eth_getTransactionByHash", hash)

	if err != nil {
		return nil, err
	} else if tx == nil {
		return nil, ethereum.NotFound
	}
	return tx, nil
}

// TxCountByAddress implements EthClient.
func (c *Clnt) TxCountByAddress(address common.Address) (hexutil.Uint64, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()
	var nonce hexutil.Uint64
	err := c.rpc.CallContext(ctxwt, &nonce, "eth_getTransactionCount", address, "latest")
	if err != nil {
		return 0, err
	}

	return nonce, err
}

// TxReceiptByHash implements EthClient.
func (c *Clnt) TxReceiptByHash(hash common.Hash) (*types.Receipt, error) {
	ctxwt, close := context.WithTimeout(context.Background(), defaultRequestTimeout.Truncate())
	defer close()

	var txReceipt *types.Receipt
	err := c.rpc.CallContext(ctxwt, &txReceipt, "eth_getTransactionReceipt", hash)
	if err != nil {
		return nil, err
	}
	if txReceipt == nil {
		return nil, ethereum.NotFound
	}
	return txReceipt, nil
}

func DiglEthClient(ctx context.Context, rpcUrl string) (EthClient, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultDialTimeout)
	defer cancel()

	bOff := retry.Exponential()
	rpcClient, err := retry.Do(ctx, defaultDialAttempts, bOff, func() (*rpc.Client, error) {
		if !IsURLAvailable(rpcUrl) {
			return nil, fmt.Errorf("address unavailable (%s)", rpcUrl)
		}

		client, err := rpc.DialContext(ctx, rpcUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to dial address (%s): %w", rpcUrl, err)
		}

		return client, nil
	})

	if err != nil {
		return nil, err
	}

	return &Clnt{rpc: NewRPC(rpcClient)}, nil
}

type Logs struct {
	Logs          []types.Log
	ToBlockHeader *types.Header
}

type RPC interface {
	Close()
	CallContext(ctx context.Context, result any, method string, args ...any) error
	BatchCallContext(ctx context.Context, b []rpc.BatchElem) error
}

type rpcClient struct {
	rpc *rpc.Client
}

func NewRPC(client *rpc.Client) RPC {
	return &rpcClient{client}
}

func (c *rpcClient) Close() {
	c.rpc.Close()
}

func (c *rpcClient) CallContext(ctx context.Context, result any, method string, args ...any) error {
	err := c.rpc.CallContext(ctx, result, method, args...)
	return err
}

func (c *rpcClient) BatchCallContext(ctx context.Context, b []rpc.BatchElem) error {
	err := c.rpc.BatchCallContext(ctx, b)
	return err
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	if number.Sign() >= 0 {
		return hexutil.EncodeBig(number)
	}
	return rpc.BlockNumber(number.Int64()).String()
}

func toFilterArg(q ethereum.FilterQuery) (interface{}, error) {
	arg := map[string]interface{}{"address": q.Addresses, "topics": q.Topics}
	if q.BlockHash != nil {
		arg["blockHash"] = *q.BlockHash
		if q.FromBlock != nil || q.ToBlock != nil {
			return nil, errors.New("cannot specify both BlockHash and FromBlock/ToBlock")
		}
	} else {
		if q.FromBlock == nil {
			arg["fromBlock"] = "0x0"
		} else {
			arg["fromBlock"] = toBlockNumArg(q.FromBlock)
		}
		arg["toBlock"] = toBlockNumArg(q.ToBlock)
	}
	return arg, nil
}

func IsURLAvailable(address string) bool {
	u, err := url.Parse(address)
	if err != nil {
		return false
	}
	addr := u.Host
	if u.Port() == "" {
		switch u.Scheme {
		case "http", "ws":
			addr += ":80"
		case "https", "wss":
			addr += ":443"
		default:
			// Fail open if we can't figure out what the port should be
			return true
		}
	}
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
