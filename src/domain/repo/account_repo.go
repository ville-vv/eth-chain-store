package repo

import (
	"github.com/ville-vv/eth-chain-store/src/infra/dao"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
)

type EthAccountRepository interface {
	//
	QueryNormalAccount(addr string) (*model.EthereumAccount, string, error)
	//
	UpdateBalanceById(tableName string, id int64, balance string) error
	//
	CreateEthAccount(normalAccount *model.EthereumAccount) error
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
	tbName, err := sel.accountDao.QueryNormalAccountInCurrentTb(addr, &info)
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
