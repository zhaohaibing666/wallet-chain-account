package ethereum

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/zhaohaibing666/wallet-chain-account/chain"
	"github.com/zhaohaibing666/wallet-chain-account/config"
	"github.com/zhaohaibing666/wallet-chain-account/rpc/account"
	common2 "github.com/zhaohaibing666/wallet-chain-account/rpc/common"
)

const ChainName = "Ethereum"

type ChainAdaptor struct {
	ethClient     EthClient
	ethDataClient *EthData
}

// BuildSignedTransaction implements chain.IChainAdaptor.
func (c *ChainAdaptor) BuildSignedTransaction(*account.SignedTransactionRequest) (*account.SignedTransactionResponse, error) {
	panic("unimplemented")
}

// ConvertAddress implements chain.IChainAdaptor.
func (c *ChainAdaptor) ConvertAddress(*account.ConvertAddressRequest) (*account.ConvertAddressResponse, error) {
	panic("unimplemented")
}

// CreateUnSignTransaction implements chain.IChainAdaptor.
func (c *ChainAdaptor) CreateUnSignTransaction(*account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error) {
	panic("unimplemented")
}

// DecodeTransaction implements chain.IChainAdaptor.
func (c *ChainAdaptor) DecodeTransaction(*account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error) {
	panic("unimplemented")
}

// GetAccount implements chain.IChainAdaptor.
func (c *ChainAdaptor) GetAccount(*account.AccountRequest) (*account.AccountResponse, error) {
	panic("unimplemented")
}

// GetBlockByHash implements chain.IChainAdaptor.
func (c *ChainAdaptor) GetBlockByHash(req *account.BlockHashRequest) (*account.BlockResponse, error) {
	block, err := c.ethClient.BlockByHash(common.HexToHash(req.Hash))
	if err != nil {
		return &account.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "blcok by hash error",
		}, nil
	}
	var txListRet []*account.BlockInfoTransactionList
	for _, v := range block.Transactions {
		bitlItem := &account.BlockInfoTransactionList{
			From:   v.From,
			To:     v.To,
			Amount: v.Value,
			Hash:   v.Hash,
		}
		txListRet = append(txListRet, bitlItem)
	}
	return &account.BlockResponse{
		Code:         common2.ReturnCode_SUCCESS,
		Msg:          "block by hash success",
		Height:       int64(block.Height),
		Hash:         block.Hash.String(),
		BaseFee:      block.BaseFee,
		Transactions: txListRet,
	}, nil
}

// GetBlockByNumber implements chain.IChainAdaptor.
func (c *ChainAdaptor) GetBlockByNumber(req *account.BlockNumberRequest) (*account.BlockResponse, error) {
	block, err := c.ethClient.BlockByNumber(big.NewInt(req.Height))
	if err != nil {
		return &account.BlockResponse{
			Code: common2.ReturnCode_ERROR,
			Msg:  "block by number error",
		}, nil
	}
	var txListRet []*account.BlockInfoTransactionList
	for _, v := range block.Transactions {
		bitlItem := &account.BlockInfoTransactionList{
			From:   v.From,
			To:     v.To,
			Hash:   v.Hash,
			Amount: v.Value,
		}
		txListRet = append(txListRet, bitlItem)
	}
	return &account.BlockResponse{
		Code:         common2.ReturnCode_SUCCESS,
		Msg:          "block by number success",
		Height:       int64(block.Height),
		Hash:         block.Hash.String(),
		BaseFee:      block.BaseFee,
		Transactions: txListRet,
	}, nil
}

// GetBlockByRange implements chain.IChainAdaptor.
func (c *ChainAdaptor) GetBlockByRange(req *account.BlockByRangeRequest) (*account.BlockByRangeResponse, error) {
	panic("unimplemented")
}

// GetBlockHeaderByHash implements chain.IChainAdaptor.
func (c *ChainAdaptor) GetBlockHeaderByHash(req *account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error) {
	panic("unimplemented")
}

// GetBlockHeaderByNumber implements chain.IChainAdaptor.
func (c *ChainAdaptor) GetBlockHeaderByNumber(*account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error) {
	panic("unimplemented")
}

// GetExtraData implements chain.IChainAdaptor.
func (c *ChainAdaptor) GetExtraData(*account.ExtraDataRequest) (*account.ExtraDataResponse, error) {
	panic("unimplemented")
}

// GetFee implements chain.IChainAdaptor.
func (c *ChainAdaptor) GetFee(*account.FeeRequest) (*account.FeeResponse, error) {
	panic("unimplemented")
}

// GetSupportChains implements chain.IChainAdaptor.
func (c *ChainAdaptor) GetSupportChains(*account.SupportChainsRequest) (*account.SupportChainsResponse, error) {
	panic("unimplemented")
}

// GetTxByAddress implements chain.IChainAdaptor.
func (c *ChainAdaptor) GetTxByAddress(*account.TxAddressRequest) (*account.TxAddressResponse, error) {
	panic("unimplemented")
}

// GetTxByHash implements chain.IChainAdaptor.
func (c *ChainAdaptor) GetTxByHash(*account.TxHashRequest) (*account.TxHashResponse, error) {
	panic("unimplemented")
}

// SendTx implements chain.IChainAdaptor.
func (c *ChainAdaptor) SendTx(*account.SendTxRequest) (*account.SendTxResponse, error) {
	panic("unimplemented")
}

// ValidAddress implements chain.IChainAdaptor.
func (c *ChainAdaptor) ValidAddress(*account.ValidAddressRequest) (*account.ValidAddressResponse, error) {
	panic("unimplemented")
}

// VerifySignedTransaction implements chain.IChainAdaptor.
func (c *ChainAdaptor) VerifySignedTransaction(*account.VerifySignedTxRequest) (*account.VerifySignedTxResponse, error) {
	panic("unimplemented")
}

func NewChainAdaptor(conf *config.Config) (chain.IChainAdaptor, error) {
	ethClient, err := DiglEthClient(context.Background(), conf.WalletNode.Eth.RPCs[0].RPCURL)
	if err != nil {
		return nil, err
	}
	ethDataClient, err := NewEthDataClient(conf.WalletNode.Eth.DataApiUrl, conf.WalletNode.Eth.DataApiKey, time.Second*15)
	if err != nil {
		return nil, err
	}
	return &ChainAdaptor{
		ethClient:     ethClient,
		ethDataClient: ethDataClient,
	}, nil
}
