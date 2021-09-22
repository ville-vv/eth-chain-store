package repo

import (
	"github.com/ville-vv/eth-chain-store/src/infra/dao"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
)

type ContractAccountRepository interface {
	QueryBindContractAccount(addr string, contractAddr string) (*model.ContractAccountBind, string, error)
	//
	CreateContractAccount(bind *model.ContractAccountBind) error
	//
	UpdateBalanceById(tableName string, id int64, balance string) error
}

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

func (sel *ContractAccountRepo) CreateContractAccount(bind *model.ContractAccountBind) error {
	return sel.accountDao.CreateContractAccount(bind)
}
