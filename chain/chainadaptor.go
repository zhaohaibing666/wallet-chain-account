package chain

import "github.com/zhaohaibing666/wallet-chain-account/rpc/account"

type IChainAdaptor interface {
	GetSupportChains(*account.SupportChainsRequest) (*account.SupportChainsResponse, error)

	ConvertAddress(*account.ConvertAddressRequest) (*account.ConvertAddressResponse, error)

	ValidAddress(*account.ValidAddressRequest) (*account.ValidAddressResponse, error)

	GetBlockByNumber(*account.BlockNumberRequest) (*account.BlockResponse, error)

	GetBlockByHash(*account.BlockHashRequest) (*account.BlockResponse, error)

	GetBlockHeaderByNumber(*account.BlockHeaderNumberRequest) (*account.BlockHeaderResponse, error)

	GetBlockHeaderByHash(*account.BlockHeaderHashRequest) (*account.BlockHeaderResponse, error)

	GetAccount(*account.AccountRequest) (*account.AccountResponse, error)

	GetFee(*account.FeeRequest) (*account.FeeResponse, error)

	SendTx(*account.SendTxRequest) (*account.SendTxResponse, error)

	GetTxByAddress(*account.TxAddressRequest) (*account.TxAddressResponse, error)

	GetTxByHash(*account.TxHashRequest) (*account.TxHashResponse, error)

	GetBlockByRange(*account.BlockByRangeRequest) (*account.BlockByRangeResponse, error)

	CreateUnSignTransaction(*account.UnSignTransactionRequest) (*account.UnSignTransactionResponse, error)

	BuildSignedTransaction(*account.SignedTransactionRequest) (*account.SignedTransactionResponse, error)

	DecodeTransaction(*account.DecodeTransactionRequest) (*account.DecodeTransactionResponse, error)

	VerifySignedTransaction(*account.VerifySignedTxRequest) (*account.VerifySignedTxResponse, error)

	GetExtraData(*account.ExtraDataRequest) (*account.ExtraDataResponse, error)
}
