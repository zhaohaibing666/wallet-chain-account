package ethereum

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/url"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
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
func (c *Clnt) BlockByHash(common.Hash) (*RpcBlock, error) {
	panic("unimplemented")
}

// BlockByNumber implements EthClient.
func (c *Clnt) BlockByNumber(*big.Int) (*RpcBlock, error) {
	panic("unimplemented")
}

// BlockHeaderByHash implements EthClient.
func (c *Clnt) BlockHeaderByHash(common.Hash) (*types.Header, error) {
	panic("unimplemented")
}

// BlockHeaderByNumber implements EthClient.
func (c *Clnt) BlockHeaderByNumber(*big.Int) (*types.Header, error) {
	panic("unimplemented")
}

// BlockHeadersByRange implements EthClient.
func (c *Clnt) BlockHeadersByRange(*big.Int, *big.Int, uint) ([]types.Header, error) {
	panic("unimplemented")
}

// Close implements EthClient.
func (c *Clnt) Close() {
	panic("unimplemented")
}

// FilterLogs implements EthClient.
func (c *Clnt) FilterLogs(filterQuery ethereum.FilterQuery, chainId uint) (Logs, error) {
	panic("unimplemented")
}

// LatestFinalizedBlockHeader implements EthClient.
func (c *Clnt) LatestFinalizedBlockHeader() (*types.Header, error) {
	panic("unimplemented")
}

// LatestSafeBlockHeader implements EthClient.
func (c *Clnt) LatestSafeBlockHeader() (*types.Header, error) {
	panic("unimplemented")
}

// SendRawTransaction implements EthClient.
func (c *Clnt) SendRawTransaction(rawTx string) error {
	panic("unimplemented")
}

// StorageHash implements EthClient.
func (c *Clnt) StorageHash(common.Address, *big.Int) (common.Hash, error) {
	panic("unimplemented")
}

// SuggestGasPrice implements EthClient.
func (c *Clnt) SuggestGasPrice() (*big.Int, error) {
	panic("unimplemented")
}

// SuggestGasTipCap implements EthClient.
func (c *Clnt) SuggestGasTipCap() (*big.Int, error) {
	panic("unimplemented")
}

// TxByHash implements EthClient.
func (c *Clnt) TxByHash(common.Hash) (*types.Transaction, error) {
	panic("unimplemented")
}

// TxCountByAddress implements EthClient.
func (c *Clnt) TxCountByAddress(common.Address) (hexutil.Uint64, error) {
	panic("unimplemented")
}

// TxReceiptByHash implements EthClient.
func (c *Clnt) TxReceiptByHash(common.Hash) (*types.Receipt, error) {
	panic("unimplemented")
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
