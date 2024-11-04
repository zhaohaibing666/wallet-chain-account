package ethereum

import (
	"context"
	"time"

	"github.com/zhaohaibing666/wallet-chain-account/chain"
	"github.com/zhaohaibing666/wallet-chain-account/config"
	"github.com/zhaohaibing666/wallet-chain-account/rpc/accout"
)

const ChainName = "Ethereum"

type ChainAdaptor struct {
	ethClient     EthClient
	ethDataClient *EthData
}

// BuildSignedTransaction implements chain.IChainAdaptor.
func (c *ChainAdaptor) BuildSignedTransaction(*accout.SignedTransactionRequest) (*accout.SignedTransactionResponse, error) {
	panic("unimplemented")
}

// ConvertAddress implements chain.IChainAdaptor.
func (c *ChainAdaptor) ConvertAddress(*accout.ConvertAddressRequest) (*accout.ConvertAddressReponse, error) {
	panic("unimplemented")
}

// CreateUnsignTransaction implements chain.IChainAdaptor.
func (c *ChainAdaptor) CreateUnsignTransaction(*accout.UnSignTransactionRequest) (*accout.UnSignTransactionResponse, error) {
	panic("unimplemented")
}

// DecodeTransaction implements chain.IChainAdaptor.
func (c *ChainAdaptor) DecodeTransaction(*accout.DecodeTransactionRequest) (*accout.DecodeTransactionReponse, error) {
	panic("unimplemented")
}

// GetAccount implements chain.IChainAdaptor.
func (c *ChainAdaptor) GetAccount(*accout.AccountRequest) (*accout.AccountResponse, error) {
	panic("unimplemented")
}

// GetBlockByHash implements chain.IChainAdaptor.
func (c *ChainAdaptor) GetBlockByHash(*accout.BlockHashRequest) (*accout.BlockReonse, error) {
	panic("unimplemented")
}

// GetBlockByNumber implements chain.IChainAdaptor.
func (c *ChainAdaptor) GetBlockByNumber(*accout.BlockNumberRequest) (*accout.BlockReonse, error) {
	panic("unimplemented")
}

// GetBlockByRange implements chain.IChainAdaptor.
func (c *ChainAdaptor) GetBlockByRange(*accout.BlockByRangeRequest) (*accout.BlockByRangeReponse, error) {
	panic("unimplemented")
}

// GetBlockHeaderByHash implements chain.IChainAdaptor.
func (c *ChainAdaptor) GetBlockHeaderByHash(*accout.BlockHeaderHashRequest) (*accout.BlockHeaderResponse, error) {
	panic("unimplemented")
}

// GetBlockHeaderByNumber implements chain.IChainAdaptor.
func (c *ChainAdaptor) GetBlockHeaderByNumber(*accout.BlockHeaderNumberRequest) (*accout.BlockHeaderResponse, error) {
	panic("unimplemented")
}

// GetExtraData implements chain.IChainAdaptor.
func (c *ChainAdaptor) GetExtraData(*accout.ExtraDataRequest) (*accout.ExtraDataReponse, error) {
	panic("unimplemented")
}

// GetFee implements chain.IChainAdaptor.
func (c *ChainAdaptor) GetFee(*accout.FeeRequest) (*accout.FeeResponse, error) {
	panic("unimplemented")
}

// GetSupportChains implements chain.IChainAdaptor.
func (c *ChainAdaptor) GetSupportChains(*accout.SupportChainsRequest) (*accout.SupportChainsReponse, error) {
	panic("unimplemented")
}

// GetTxByAddress implements chain.IChainAdaptor.
func (c *ChainAdaptor) GetTxByAddress(*accout.TxAddressRequest) (*accout.TxAddressResponse, error) {
	panic("unimplemented")
}

// GetTxByHash implements chain.IChainAdaptor.
func (c *ChainAdaptor) GetTxByHash(*accout.TxHashRequest) (*accout.TxHashResponse, error) {
	panic("unimplemented")
}

// SendTx implements chain.IChainAdaptor.
func (c *ChainAdaptor) SendTx(*accout.SendTxRequest) (*accout.SendTxResponse, error) {
	panic("unimplemented")
}

// ValidAdress implements chain.IChainAdaptor.
func (c *ChainAdaptor) ValidAdress(*accout.ValidAddressRequest) (*accout.ValidAddressResponse, error) {
	panic("unimplemented")
}

// VerifySignedTransaction implements chain.IChainAdaptor.
func (c *ChainAdaptor) VerifySignedTransaction(*accout.VerifySignedTxRequest) (*accout.VerifySignedTxResponse, error) {
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
