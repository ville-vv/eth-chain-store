package repo

import "github.com/ville-vv/eth-chain-store/src/infra/model"

type ContractAccountRepository interface {
	IsAccountExist(addr string) (bool, error)
	//
	CreateAccount(account *model.EthereumAccount) error
	//
	UpdateAccountByAddr(addr string, updateInfo map[string]interface{}) error
}
