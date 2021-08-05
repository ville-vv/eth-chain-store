package ethm

import (
	"github.com/pkg/errors"
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/eth-chain-store/src/infra/mqp"
	"github.com/ville-vv/vilgo/vlog"
)

type Account struct {
	Address    string
	IsContract bool
	Balance    string
}

// 以太坊账户
type AccountManager struct {
	ethCli         ethrpc.EthRPC
	contractActMng *contractAccountManager
	normalActMng   *normalAccountManager
}

func NewAccountManager(ethCli ethrpc.EthRPC, contractAccountRepo *repo.ContractAccountRepo, normalAccountRepo *repo.NormalAccountRepo) *AccountManager {
	ac := &AccountManager{
		ethCli: ethCli,
		//contractMng: contractMng,
		contractActMng: &contractAccountManager{
			ethCli:      ethCli,
			accountRepo: contractAccountRepo,
		},
		normalActMng: &normalAccountManager{
			ethCli:      ethCli,
			accountRepo: normalAccountRepo,
		},
	}
	return ac
}

// TxWrite 写入代币合约与钱包地址绑定，在一笔交易中，外部的普通交易 to 地址为合约地址，from 地址为钱包地址，
// 如果交易类型是合约代币交易，那么合约地址为 contractAddress
// 内部交易无法确定,就当做是普通账户写入
func (sel *AccountManager) TxWrite(txData *model.TransactionData) error {
	vlog.DEBUG("tx writer to account manager %v %s ", txData.IsLatest(), txData.Hash)
	if txData.IsContractToken {
		return sel.contractActMng.UpdateAccount(txData)
	}
	return sel.normalActMng.UpdateAccount(txData)
}

func (sel *AccountManager) ID() string {
	return "AccountManager"
}
func (sel *AccountManager) Process(msg *mqp.Message) error {
	txData := &model.TransactionData{}
	err := msg.UnMarshalFromBody(txData)
	if err != nil {
		return err
	}

	return sel.TxWrite(txData)
}

// contractAccountUpdater 处理合约账户
func (sel *AccountManager) contractAccountUpdater(txData *model.TransactionData) error {

	//accountRepo.FromBalance = balance
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
	accountRepo *repo.ContractAccountRepo
}

// UpdateAccount 合约代币交易账户信息写入, contractAddress 是合约地址
func (sel *contractAccountManager) UpdateAccount(txData *model.TransactionData) error {
	if err := sel.writeAccount(txData.From, txData.ContractAddress, txData.IsLatest()); err != nil {
		return errors.Wrap(err, "write contract account from address")
	}
	if err := sel.writeAccount(txData.To, txData.ContractAddress, txData.IsLatest()); err != nil {
		return errors.Wrap(err, "write contract account to address")
	}
	return nil
}

func (sel *contractAccountManager) writeAccount(accountAddr string, contractAddress string, isLatest bool) error {
	balance, err := sel.ethCli.GetContractBalance(contractAddress, accountAddr)
	if err != nil {
		return err
	}

	updated, err := sel.accountRepo.UpdateBalance(accountAddr, contractAddress, balance, isLatest)
	if err != nil {
		vlog.ERROR("writeAccount update account balance failed addr:%s contract:%s error:%s", accountAddr, contractAddress, err.Error())
		return err
	}
	if updated {
		return nil
	}

	symbol, err := sel.ethCli.GetContractSymbol(contractAddress)
	if err != nil {
		vlog.ERROR("writeAccount get contract symbol failed addr:%s contract:%s", accountAddr, contractAddress)
		return err
	}
	// 如果该账户不存在
	return sel.accountRepo.CreateEthAccount(&model.ContractAccountBind{
		Address:         accountAddr,
		ContractAddress: contractAddress,
		Symbol:          symbol,
		Balance:         balance,
	})
}

// normalAccountManager 以太坊账户管理
type normalAccountManager struct {
	ethCli      ethrpc.EthRPC
	accountRepo *repo.NormalAccountRepo
}

// 以太坊正常的交易账户写入，这里就不判断该账户是不是合约账户了直接写入 from 和 to
func (sel *normalAccountManager) UpdateAccount(txData *model.TransactionData) error {
	if err := sel.writeAccount(txData.From, txData.TimeStamp, txData.Hash, txData.IsLatest()); err != nil {
		return errors.Wrap(err, "write normal account from address")
	}
	if err := sel.writeAccount(txData.To, txData.TimeStamp, txData.Hash, txData.IsLatest()); err != nil {
		return errors.Wrap(err, "write normal account to address")
	}
	return nil
}
func (sel *normalAccountManager) writeAccount(accountAddr string, timeStamp string, hash string, isLatest bool) error {
	// 获取的余额
	balance, err := sel.ethCli.GetBalance(accountAddr)
	if err != nil {
		return err
	}

	updated, err := sel.accountRepo.UpdateBalance(accountAddr, balance, isLatest)
	if err != nil {
		return err
	}
	if updated {
		// 账户存在，已经更新
		return nil
	}
	// 创建一个以太坊账户
	return sel.accountRepo.CreateEthAccount(&model.EthereumAccount{
		Address:     accountAddr,
		FirstTxTime: timeStamp,
		FirstTxHash: hash,
		Balance:     balance,
	})
}
