package ethm

import "github.com/ville-vv/eth-chain-store/src/infra/ethrpc"

type EthRpcExecutor struct {
	ethrpc.EthRPC
}

func NewEthRpcExecutor(endPoint string) *EthRpcExecutor {
	ethRPC := ethrpc.NewClient(endPoint)

	return &EthRpcExecutor{EthRPC: ethRPC}
}
