package ethm

import (
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
)

type Account struct {
	Address    string
	IsContract bool
	Balance    string
}

// 以太坊账户
type AccountManager struct {
	ethCli         ethrpc.EthRPC
	accountRepo    repo.NormalAccountRepo
	contractMng    ContractManager
	contractActMng *contractAccountManager
	normalActMng   *normalAccountManager
}

func (sel *AccountManager) TxWrite(txData *model.TransactionData) error {
	if txData.IsContract {
		return sel.contractActMng.UpdateAccount(txData)
	}
	return sel.normalActMng.UpdateAccount(txData)
}

// contractAccountUpdater 处理合约账户
func (sel *AccountManager) contractAccountUpdater(txData *model.TransactionData) error {

	//accountRepo.Balance = balance
	//// 在数据库中是否存在该合约账户
	//if sel.accountRepo.IsContractAccountExist() {
	//	// 如果存在就直接更新余额
	//	if err = sel.accountRepo.UpdateContract(act); err != nil {
	//		return err
	//	}
	//}
	//
	//// 如果不存在就需要获取该合约相关的东西
	//totalSupply, err := sel.ethCli.GetContractTotalSupply(act.Address)
	//if err != nil {
	//	return err
	//}

	//sel.ethCl
	return nil
}

// contractAccountManager 以太坊合约账户管理
type contractAccountManager struct {
	ethCli      ethrpc.EthRPC
	accountRepo repo.ContractAccountRepo
}

// UpdateAccount 更新合约账户信息，如果该用户存在就更新余额信息，如果该用户不存在就创建一条记录
// 以交易记录的 to 地址来进行判断
func (sel *contractAccountManager) UpdateAccount(txData *model.TransactionData) error {
	accountAddr := txData.To
	balance, err := sel.ethCli.GetContractBalance(txData.ContractAddress, accountAddr)
	if err != nil {
		return err
	}
	ok, err := sel.accountRepo.IsAccountExist(accountAddr, txData.ContractAddress)
	if err != nil {
		return err
	}
	if ok {
		// 该账户已经存在
		return sel.accountRepo.UpdateBalance(accountAddr, txData.ContractAddress, balance)
	}
	symbol, err := sel.ethCli.GetContractSymbol(txData.ContractAddress)
	if err != nil {
		return err
	}
	// 如果该账户不存在
	return sel.accountRepo.CreateEthAccount(&model.ContractAccountBind{
		Address:         accountAddr,
		ContractAddress: txData.ContractAddress,
		Symbol:          symbol,
		Balance:         balance,
	})
}

// normalAccountManager 以太坊账户管理
type normalAccountManager struct {
	ethCli      ethrpc.EthRPC
	accountRepo repo.NormalAccountRepo
}

func (sel *normalAccountManager) UpdateAccount(txData *model.TransactionData) error {
	var addr = txData.To
	ok, err := sel.accountRepo.IsAccountExist(addr)
	if err != nil {
		return err
	}
	// 获取的余额
	balance, err := sel.ethCli.GetBalance(addr)
	if err != nil {
		return err
	}
	if ok {

		return sel.accountRepo.UpdateBalance(addr, balance)
	}
	// 创建一个以太坊账户
	return sel.accountRepo.CreateEthAccount(&model.EthereumAccount{
		Address:     txData.To,
		FirstTxTime: txData.TimeStamp,
		FirstTxHash: txData.Hash,
		Balance:     balance,
	})
}
