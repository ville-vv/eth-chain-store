package ethrpc

type EthRPC interface {
	// GetBlockNumber 获取最新区块号
	GetBlockNumber() (uint64, error)

	// GetBlockByLatest
	GetBlock() (*EthBlock, error)

	GetBlockByNumber(blockNumber int64) (*EthBlock, error)

	GetContractTotalSupply(contract string) (string, error)

	GetBalance(addr string) (string, error)

	GetBalanceByBlockNumber(addr string, blockNumber int64) (string, error)

	GetContractBalance(contract, addr string) (string, error)

	GetTransactionReceipt(hash string) (*EthTransactionReceipt, error)

	GetContractBalanceByBlockNumber(contract, addr string, blockNumber int64) (string, error)

	GetCode(addr string) (string, error)

	GetContractSymbol(contract string) (string, error)

	GetContractName(contract string) (string, error)

	GetContractDecimals(contract string) (string, error)
}
