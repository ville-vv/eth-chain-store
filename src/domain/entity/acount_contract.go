package entity

import (
	"github.com/ville-vv/eth-chain-store/src/common/utils"
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"time"
)

type AccountContract struct {
	model.EthereumAccount
	rpcCli      ethrpc.EthRPC
	contract    *model.ContractContent
	accountRepo repo.ContractAccountRepository
}

func NewAccountContract(accountRepo repo.ContractAccountRepository) *AccountContract {
	return &AccountContract{accountRepo: accountRepo}
}

func (sel *AccountContract) SetAddress(addr string) {
	sel.Address = addr
	return
}

func (sel *AccountContract) SetFirstTxTime(tmStr string) {
	sel.FirstTxTime = tmStr
}

func (sel *AccountContract) SetContractContent(contractContent *model.ContractContent) {
	sel.contract = contractContent
}

func (sel *AccountContract) CreateRecord() error {
	//contractAccount, tbName, err := sel.accountRepo.QueryBindContractAccount(sel.Address, sel.contract.Address)
	//if err != nil {
	//	return err
	//}
	//if contractAccount.ID > 0 {
	//
	//	return nil
	//}
	//
	//if err := sel.SetBalanceFromRpc(); err != nil {
	//	return err
	//}
	//
	//return sel.accountRepo.CreateContractAccount(&model.ContractAccountBind{
	//	TxTime:          sel.GetTxTime(),
	//	Address:         sel.Address,
	//	ContractAddress: sel.contract.Address,
	//	Symbol:          sel.contract.Symbol,
	//	Balance:         sel.Balance,
	//})
	return nil
}

func (sel *AccountContract) GetTxTime() *time.Time {
	timeFm, _ := utils.ParseLocal(sel.FirstTxTime)
	return &timeFm
}

func (sel *AccountContract) SetBalanceFromRpc() error {
	balance, err := sel.rpcCli.GetContractBalance(sel.contract.Address, sel.Address)
	if err != nil {
		return err
	}
	sel.Balance = balance
	return nil
}

func (sel *AccountContract) UpdateBalance() error {
	return nil
}
