package ethm

import (
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"strconv"
)

type Erc20Manager interface {
	IsErc20(addr string) bool
}

type Erc20Contract struct {
	Address     string // 合约地址
	Name        string // 名字
	Symbol      string // 标识
	TotalSupply string // 发行量
	Decimal     int
	Balance     string
}

type ContractManager struct {
	rpcCli       ethrpc.EthRPC
	contractRepo repo.ContractRepo
	erc20Mng     Erc20Manager
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

	balance, err := sel.rpcCli.GetContractBalance(contractAddr, contractAddr)
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

	return &Erc20Contract{
		Address:     contractAddr,
		Name:        name,
		Symbol:      symbol,
		TotalSupply: supply,
		Decimal:     int(decimalInt),
		Balance:     balance,
	}, nil
}

func (sel *ContractManager) TxWrite(txData *model.TransactionData) (err error) {
	// 查询是否已经存在记录
	if sel.contractRepo.IsContractExist(txData.ContractAddress) {
		return nil
	}
	var erc20ContractInfo = &Erc20Contract{}
	// 如果不存在就创建
	if txData.IsErc20 {
		// erc20 合约
		erc20ContractInfo, err = sel.GetErc20ContractInfo(txData.ContractAddress)
		if err != nil {
			return err
		}
	}
	return sel.contractRepo.CreateContract(&model.ContractAccount{
		Symbol:      erc20ContractInfo.Symbol,
		Address:     txData.ContractAddress,
		PublishTime: txData.TimeStamp,
		IsErc20:     txData.IsErc20,
		TotalSupply: erc20ContractInfo.TotalSupply,
	})
}
