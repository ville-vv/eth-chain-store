package dao

import "github.com/ville-vv/eth-chain-store/src/infra/model"

type EthereumDao struct {
}

// QueryContractInfo 查询合约信息
func (sel *EthereumDao) QueryContractInfo(addr string, contractInfo *model.TbContractAddress) error {
	return nil
}

// CreateContractRecord 创建合约记录
func (sel *EthereumDao) CreateContractRecord(contractCtx *model.ContractContent) error {
	return nil
}

// QueryNormalAccount 查询普通的以太坊账户地址信息
func (sel *EthereumDao) QueryNormalAccount(addr string, info *model.EthereumAccount) error {
	return nil
}

// QueryBindContractAccount 查询合约绑定账户信息
func (sel *EthereumDao) QueryBindContractAccount(addr, contractAddr string, bindInfo *model.ContractAccountBind) error {
	return nil
}

// UpdateNativeBalance 更新以太坊
func (sel *EthereumDao) UpdateNormalAccountBalance(addr string, balance string) error {
	return nil
}

func (sel *EthereumDao) CreateNormalAccount(normalAccount *model.EthereumAccount) error {
	return nil
}
