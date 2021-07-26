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

func (sel *Erc20Contract) IsErc20() bool {
	return sel.Name != "" && sel.Symbol != "" && sel.TotalSupply != "" && sel.Decimal != 0
}

type ContractManager struct {
	rpcCli       ethrpc.EthRPC
	contractRepo *repo.ContractRepo
}

func NewContractManager(rpcCli ethrpc.EthRPC, contractRepo *repo.ContractRepo) *ContractManager {
	return &ContractManager{rpcCli: rpcCli, contractRepo: contractRepo}
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

	var decimalInt int64
	if decimal != "0x" {
		decimalInt, err = strconv.ParseInt(decimal, 0, 64)
		if err != nil {
			return nil, errors.Wrap(err, "parse int contract decimal value "+decimal)
		}
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

// TxWrite 合约信息写入，一笔交易存在两个地址，一个是from 地址，一个是 to 地址，两个地址都有可能是合约地址，
// 如果是  token transfer 交易，那么 to 地址一定是合约地址
func (sel *ContractManager) TxWrite(txData *model.TransactionData) (err error) {
	if txData.IsContractToken {
		return sel.writeTokenContractInfo(txData.ContractAddress, txData.TimeStamp)
	}

	// 检查 from 地址
	if err = sel.writeOtherContractInfo(txData.From, txData.TimeStamp); err != nil {
		return nil
	}

	// 检查 to 地址
	if err = sel.writeOtherContractInfo(txData.To, txData.TimeStamp); err != nil {
		return nil
	}

	return nil
}

// writeTokenContractInfo 代币合约信息
func (sel *ContractManager) writeTokenContractInfo(addr string, timeStamp string) (err error) {
	// 如果存在合约地址，也要到主链中判断该地址是不是合约地址
	codeData, err := sel.rpcCli.GetCode(addr)
	if err != nil {
		return err
	}
	if codeData == "0x" || codeData == "" {
		// 不是合约地址直接返回
		return nil
	}

	// 查询是否已经存在记录
	if sel.contractRepo.IsContractExist(addr) {
		return nil
	}

	// erc20 合约
	erc20ContractInfo, err := sel.GetErc20ContractInfo(addr)
	if err != nil {
		return errors.Wrap(err, "get erc20 contract info")
	}
	// 如果不存在就创建
	return sel.contractRepo.CreateContract(&model.ContractContent{
		Symbol:      erc20ContractInfo.Symbol,
		Address:     addr,
		PublishTime: timeStamp,
		IsErc20:     true,
		TotalSupply: erc20ContractInfo.TotalSupply,
	})
}

func (sel *ContractManager) writeOtherContractInfo(addr string, timeStamp string) (err error) {
	// 如果存在合约地址，也要到主链中判断该地址是不是合约地址
	codeData, err := sel.rpcCli.GetCode(addr)
	if err != nil {
		return err
	}
	if codeData == "0x" || codeData == "" {
		// 不是合约地址直接返回
		return nil
	}
	// 查询是否已经存在记录
	if sel.contractRepo.IsContractExist(addr) {
		return nil
	}
	// erc20 合约
	erc20ContractInfo, err := sel.GetErc20ContractInfo(addr)
	if err != nil {
		return errors.Wrap(err, "get erc20 contract info")
	}

	return sel.contractRepo.CreateContract(&model.ContractContent{
		Symbol:      erc20ContractInfo.Symbol,
		Address:     addr,
		PublishTime: timeStamp,
		IsErc20:     erc20ContractInfo.IsErc20(),
		TotalSupply: erc20ContractInfo.TotalSupply,
	})
}
