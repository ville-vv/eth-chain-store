package repo

import (
	"github.com/ville-vv/eth-chain-store/src/infra/dao"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/vilgo/vlog"
)

type EthAccountRepository interface {
	IsAccountExist(addr string) (bool, error)
	//
	CreateAccount(account *model.EthereumAccount) error
	//
	UpdateAccountByAddr(addr string, updateInfo map[string]interface{}) error
}

//===============================================================================================
// 普通账户处理
type NormalAccountRepo struct {
	accountDao *dao.EthereumDao
}

func NewNormalAccountRepo(accountDao *dao.EthereumDao) *NormalAccountRepo {
	return &NormalAccountRepo{accountDao: accountDao}
}

func (sel *NormalAccountRepo) QueryNormalAccount(addr string) (*model.EthereumAccount, string, error) {
	var info model.EthereumAccount
	tbName, err := sel.accountDao.QueryNormalAccount(addr, &info)
	if err != nil {
		return nil, "", err
	}
	return &info, tbName, nil
}

// UpdateBalance 更新余额 数据存在就返回 true, nil
func (sel *NormalAccountRepo) UpdateBalanceById(tableName string, id int64, balance string) error {
	return sel.accountDao.UpdateNormalAccountBalanceById(tableName, id, balance)
}

// CreateEthAccount 创建地址账户
func (sel *NormalAccountRepo) CreateEthAccount(normalAccount *model.EthereumAccount) error {
	return sel.accountDao.CreateNormalAccount(normalAccount)
}

//===============================================================================================
// 合约账户处理
type ContractAccountRepo struct {
	accountDao *dao.EthereumDao
}

func NewContractAccountRepo(accountDao *dao.EthereumDao) *ContractAccountRepo {
	return &ContractAccountRepo{accountDao: accountDao}
}

func (sel *ContractAccountRepo) QueryBindContractAccount(addr string, contractAddr string) (*model.ContractAccountBind, string, error) {
	var info model.ContractAccountBind
	tbName, err := sel.accountDao.QueryBindContractAccount(addr, contractAddr, &info)
	if err != nil {
		return nil, "", err
	}
	return &info, tbName, nil
}

// UpdateNative
func (sel *ContractAccountRepo) UpdateBalanceById(tableName string, id int64, balance string) error {
	return sel.accountDao.UpdateContractAccountBalanceById(tableName, id, balance)
}
func (sel *ContractAccountRepo) CreateEthAccount(bind *model.ContractAccountBind) error {
	vlog.DEBUG("create account bind information of address:%s contract:%s", bind.Address, bind.ContractAddress)
	return sel.accountDao.CreateContractAccount(bind)
}
