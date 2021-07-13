package ethrpc

type EthRpcClient interface {
	// GetBlockNumber 获取最新区块号
	GetBlockNumber() (uint64, error)
	// GetBlockByLatest
	GetBlock() (*RpcBlock, error)
	GetBlockByNumber(blockNumber int64) (*RpcBlock, error)
	GetContractTotalSupply(contract string) (string, error)
	GetBalance(addr string) (string, error)
	GetBalanceByBlockNumber(addr string, blockNumber int64) (string, error)
	GetContractBalance(contract, addr string) (string, error)

	GetTransactionReceipt(hash string) (*RpcTransactionReceipt, error)
	GetContractBalanceByBlockNumber(contract, addr string, blockNumber int64) (string, error)
}
