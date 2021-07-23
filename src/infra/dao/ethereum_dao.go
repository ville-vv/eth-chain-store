package dao

import (
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
)

type EthereumDao struct {
	db DB
}

func NewEthereumDao(db DB) *EthereumDao {
	return &EthereumDao{}
}

// QueryContractInfo 查询合约信息
func (sel *EthereumDao) QueryContractInfo(addr string, contractInfo *model.TbContractAddress) error {
	return nil
}

// CreateContractRecord 创建合约记录
func (sel *EthereumDao) CreateContractRecord(contractCtx *model.ContractContent) error {
	fmt.Println("创建合约记录", contractCtx)
	return nil
}

// QueryBindContractAccount 查询合约绑定账户信息
func (sel *EthereumDao) QueryBindContractAccount(addr, contractAddr string, bindInfo *model.ContractAccountBind) error {
	return nil
}

// CreateContractAccount 普通钱包地址与合约地址绑定
func (sel *EthereumDao) CreateContractAccount(contractAccount *model.ContractAccountBind) error {
	fmt.Println("绑定钱包地址与合约地址", contractAccount)
	return nil
}

func (sel *EthereumDao) UpdateContractAccountBalance(addr string, contractAddr string, balance string) error {
	return nil
}

//========================================================================================================
// QueryNormalAccount 查询普通的以太坊账户地址信息
func (sel *EthereumDao) QueryNormalAccount(addr string, info *model.EthereumAccount) error {
	return nil
}

// UpdateNativeBalance 更新以太坊
func (sel *EthereumDao) UpdateNormalAccountBalance(addr string, balance string) error {
	return nil
}

// CreateNormalAccount 创建以太坊账户
func (sel *EthereumDao) CreateNormalAccount(normalAccount *model.EthereumAccount) error {
	fmt.Println("创建普通账户记录", normalAccount)
	return nil
}
