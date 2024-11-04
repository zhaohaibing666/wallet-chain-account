package ethereum

import (
	"time"

	"github.com/dapplink-labs/chain-explorer-api/explorer/etherscan"
)

type EthData struct {
	EthDataCli *etherscan.ChainExplorerAdaptor
}

func NewEthDataClient(baseUrl, apiKey string, timeout time.Duration) (*EthData, error) {
	etherscanCli, error := etherscan.NewChainExplorerAdaptor(apiKey, baseUrl, false, time.Duration(timeout))
	if error != nil {
		return nil, error
	}

	return &EthData{
		EthDataCli: etherscanCli,
	}, nil
}
