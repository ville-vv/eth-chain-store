package repo

import (
	"github.com/ville-vv/eth-chain-store/src/infra/dao"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/vilgo/vlog"
)

type NormalAccountRepo struct {
	accountDao *dao.EthereumDao
}

func (sel *NormalAccountRepo) IsAccountExist(addr string) (bool, error) {
	var info model.EthereumAccount
	err := sel.accountDao.QueryNormalAccount(addr, &info)
	if err != nil {
		vlog.ERROR("")
		return false, err
	}
	if info.ID > 0 {
		return true, nil
	}
	return false, nil
}

// UpdateContract
func (sel *NormalAccountRepo) UpdateContract(*model.AccountContent) error {
	return nil
}

// UpdateBalance 更新余额
func (sel *NormalAccountRepo) UpdateBalance(addr string, balance string) error {
	return sel.accountDao.UpdateNormalAccountBalance(addr, balance)
}

// CreateEthAccount 创建地址账户
func (sel *NormalAccountRepo) CreateEthAccount(normalAccount *model.EthereumAccount) error {
	return sel.accountDao.CreateNormalAccount(normalAccount)
}

type ContractAccountRepo struct {
	accountDao *dao.EthereumDao
}

func (sel *ContractAccountRepo) IsAccountExist(addr string, contractAddr string) (bool, error) {
	var info model.ContractAccountBind
	err := sel.accountDao.QueryBindContractAccount(addr, contractAddr, &info)
	if err != nil {
		vlog.ERROR("")
		return false, err
	}
	if info.ID > 0 {
		return true, nil
	}
	return false, nil
}

// UpdateNative
func (sel *ContractAccountRepo) UpdateBalance(addr, contractAddr, balance string) error {
	return nil
}
func (sel *ContractAccountRepo) CreateEthAccount(bind *model.ContractAccountBind) error {
	return nil
}
