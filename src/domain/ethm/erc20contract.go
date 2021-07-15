package ethm

import (
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"strconv"
)

type Erc20Contract struct {
	Address     string // 合约地址
	Name        string // 名字
	Symbol      string // 标识
	TotalSupply string // 发行量
	Decimal     int
}

type ContractManager struct {
	rpcCli ethrpc.EthRPC
}

func (sel *ContractManager) GetErc20ContractInfo(contractAddr string) (*Erc20Contract, error) {
	supply, err := sel.rpcCli.GetContractTotalSupply(contractAddr)
	if err != nil {
		return nil, err
	}

	name, err := sel.rpcCli.GetContractName(contractAddr)
	if err != nil {
		return nil, err
	}

	symbol, err := sel.rpcCli.GetContractSymbol(contractAddr)
	if err != nil {
		return nil, err
	}

	decimal, err := sel.rpcCli.GetContractDecimals(contractAddr)
	if err != nil {
		return nil, err
	}

	decimalInt, err := strconv.ParseInt(decimal, 0, 64)
	if err != nil {
		return nil, err
	}

	ctt := &Erc20Contract{
		Address:     contractAddr,
		Name:        name,
		Symbol:      symbol,
		TotalSupply: supply,
		Decimal:     int(decimalInt),
	}
	return ctt, nil
}

func (sel *ContractManager) TxWrite(txData *model.TransactionData) error {
	return nil
}
