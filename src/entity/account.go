package entity

import (
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/eth-chain-store/src/repo"
)

type Account struct {
	Address    string
	IsContract bool
	Balance    string
}

// 以太坊账户
type AccountManager struct {
	ethCli      ethrpc.EthRpcClient
	accountRepo repo.AccountRepo
}

func (sel *AccountManager) CreateOrUpdateAccount(account model.NativeAccount) error {
	balance, err := sel.ethCli.GetBalance(account.Address)
	if err != nil {
		return err
	}
	account.Balance = balance
	return nil
}

func (sel *AccountManager) CreateOrUpdateContractAccount(account *model.ContractAccount) error {
	return nil
}
