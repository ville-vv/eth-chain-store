package ethm

import (
	"github.com/pkg/errors"
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
}

func NewContractManager(rpcCli ethrpc.EthRPC) *ContractManager {
	return &ContractManager{rpcCli: rpcCli}
}

// GetErc20ContractInfo ERC20 协议的合约有固定的合约接口来获取合约的基本信息
// return *Erc20Contract 合约的基本信息
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
		return nil, errors.Wrap(err, "parse int contract decimal value "+decimal)
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
	if txData.ContractAddress == "" {
		return nil
	}
	codeData, err := sel.rpcCli.GetCode(txData.ContractAddress)
	if err != nil {
		return err
	}
	if codeData == "0x" || codeData == "" {
		return nil
	}

	// 查询是否已经存在记录
	if sel.contractRepo.IsContractExist(txData.ContractAddress) {
		return nil
	}
	var erc20ContractInfo = &Erc20Contract{}
	var isErc20 bool
	if txData.TxType == model.TxTokenTransfer {
		// erc20 合约
		erc20ContractInfo, err = sel.GetErc20ContractInfo(txData.ContractAddress)
		if err != nil {
			return err
		}
		if erc20ContractInfo.Symbol != "" && erc20ContractInfo.Name != "" {
			isErc20 = true
		}
	}
	// 如果不存在就创建
	return sel.contractRepo.CreateContract(&model.ContractContent{
		Symbol:      erc20ContractInfo.Symbol,
		Address:     txData.ContractAddress,
		PublishTime: txData.TimeStamp,
		IsErc20:     isErc20,
		TotalSupply: erc20ContractInfo.TotalSupply,
	})
}
