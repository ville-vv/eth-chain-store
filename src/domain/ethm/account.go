package ethm

import (
	"github.com/pkg/errors"
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/cache"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/vilgo/vlog"
	"strconv"
	"strings"
	"sync"
	"time"
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
			ethCli:        ethCli,
			accountRepo:   contractAccountRepo,
			haveWriteList: NewRingStrListV2(),
			contractCache: cache.NewRingCache(),
		},
		normalActMng: &normalAccountManager{
			ethCli:        ethCli,
			accountRepo:   normalAccountRepo,
			haveWriteList: NewRingStrListV2(),
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
	//return nil
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
	ethCli        ethrpc.EthRPC
	accountRepo   *repo.ContractAccountRepo
	haveWriteList *RingStrListV2
	contractCache *cache.RingCache
	sync.Mutex
}

// UpdateAccount 合约代币交易账户信息写入, contractAddress 是合约地址
func (sel *contractAccountManager) UpdateAccount(txData *model.TransactionData) error {
	//if err := sel.writeAccount(txData.From, txData.ContractAddress, txData.IsLatest()); err != nil {
	//	return errors.Wrap(err, "write contract account from address")
	//}
	if err := sel.writeAccount(txData.To, txData.ContractAddress, txData.TimeStamp, txData.IsLatest()); err != nil {
		return errors.Wrap(err, "write contract account to address")
	}
	return nil
}

func (sel *contractAccountManager) writeAccount(accountAddr string, contractAddress string, timestamp string, isLatest bool) error {
	var exist bool
	var err error
	var tableName string
	var bindInfo *model.ContractAccountBind
	var unique = accountAddr + contractAddress
	sel.Lock()
	defer sel.Unlock()
	if !sel.haveWriteList.Exist(unique) {
		sel.haveWriteList.Set(unique)
		// 如果缓存中没有，就查一下数据库
		bindInfo, tableName, err = sel.accountRepo.QueryBindContractAccount(contractAddress, accountAddr)
		if err != nil {
			sel.haveWriteList.Del(unique)
			return err
		}
		if bindInfo.ID <= 0 {
			err = sel.createAccount(accountAddr, contractAddress, timestamp)
			if err != nil {
				sel.haveWriteList.Del(unique)
				return err
			}
			return nil
		}
		exist = true
	}

	if isLatest {
		if !exist {
			// 如果缓存中没有，就查一下数据库
			bindInfo, tableName, err = sel.accountRepo.QueryBindContractAccount(contractAddress, accountAddr)
			if err != nil {
				return err
			}
			if bindInfo.ID <= 0 {
				return nil
			}
		}

		balance, err := sel.ethCli.GetContractBalance(contractAddress, accountAddr)
		if err != nil {
			if err.Error() != "execution reverted" {
				vlog.ERROR("writeAccount get contract balance update failed addr:%s contract:%s error:%s", accountAddr, contractAddress, err.Error())
				return err
			}
		}

		if err = sel.accountRepo.UpdateBalanceById(tableName, bindInfo.ID, balance); err != nil {
			vlog.ERROR("writeAccount update account balance failed addr:%s contract:%s error:%s", accountAddr, contractAddress, err.Error())
			return err
		}
	}
	return nil
}

func (sel *contractAccountManager) createAccount(accountAddr string, contractAddress string, timestamp string) error {
	balance, err := sel.ethCli.GetContractBalance(contractAddress, accountAddr)
	if err != nil {
		vlog.WARN("create contract address balance failed addr:%s contract:%s error:%s", accountAddr, contractAddress, err.Error())
		if strings.Contains(err.Error(), "invalid opcode") {
			return nil
		}
		if err.Error() != "execution reverted" {
			//vlog.ERROR("writeAccount.createAccount get contract balance failed addr:%s contract:%s error:%s", accountAddr, contractAddress, err.Error())
			return err
		}
		return nil
	}
	//symbol, err := sel.getContractSymbol(contractAddress)
	//if err != nil {
	//	vlog.ERROR("writeAccount get contract symbol failed contract:%s error:%s", contractAddress, err.Error())
	//	return err
	//}

	unixTm, _ := strconv.ParseInt(timestamp, 10, 64)
	timeFm := time.Unix(unixTm, 0)

	err = sel.accountRepo.CreateEthAccount(&model.ContractAccountBind{
		Address:         accountAddr,
		ContractAddress: contractAddress,
		Timestamp:       &timeFm,
		Balance:         balance,
	})
	if err != nil {
		return err
	}
	return nil
}

func (sel *contractAccountManager) getContractSymbol(contractAddress string) (string, error) {
	var symbol string
	var err error
	val, ok := sel.contractCache.Get(contractAddress)
	if ok {
		symbol, _ = val.(string)
	} else {
		symbol, err = sel.ethCli.GetContractSymbol(contractAddress)
		if err != nil {
			if err.Error() != "execution reverted" {
				return "", err
			}
			return "", nil
		}
	}
	return symbol, err
}

// normalAccountManager 以太坊账户管理
type normalAccountManager struct {
	ethCli        ethrpc.EthRPC
	accountRepo   *repo.NormalAccountRepo
	haveWriteList *RingStrListV2
	lock          sync.Mutex
}

// 以太坊正常的交易账户写入，这里就不判断该账户是不是合约账户了直接写入 from 和 to
func (sel *normalAccountManager) UpdateAccount(txData *model.TransactionData) error {
	//if err := sel.writeAccount(txData.From, txData.TimeStamp, txData.Hash, txData.IsLatest()); err != nil {
	//	return errors.Wrap(err, "write normal account from address")
	//}
	if err := sel.writeAccount(txData.To, txData.TimeStamp, txData.Hash, txData.IsLatest()); err != nil {
		return errors.Wrap(err, "write normal account to address")
	}
	return nil
}

func (sel *normalAccountManager) writeAccount(accountAddr string, timeStamp string, hash string, isLatest bool) error {
	var exist bool
	var err error
	var tableName string
	var bindInfo *model.EthereumAccount
	sel.lock.Lock()
	defer sel.lock.Unlock()

	if !sel.haveWriteList.Exist(accountAddr) {
		sel.haveWriteList.Set(accountAddr)
		bindInfo, tableName, err = sel.accountRepo.QueryNormalAccount(accountAddr)
		if err != nil {
			sel.haveWriteList.Del(accountAddr)
			return err
		}
		if bindInfo.ID <= 0 {
			// 获取的余额
			balance, err := sel.ethCli.GetBalance(accountAddr)
			if err != nil {
				sel.haveWriteList.Del(accountAddr)
				return err
			}
			err = sel.accountRepo.CreateEthAccount(&model.EthereumAccount{
				Address:     accountAddr,
				FirstTxTime: timeStamp,
				FirstTxHash: hash,
				Balance:     balance,
			})
			if err != nil {
				sel.haveWriteList.Del(accountAddr)
				return err
			}
			return nil
		}
		exist = true
	}

	if isLatest {
		if !exist {
			bindInfo, tableName, err = sel.accountRepo.QueryNormalAccount(accountAddr)
			if err != nil {
				return err
			}
			if bindInfo.ID <= 0 {
				return nil
			}
		}
		// 获取的余额
		balance, err := sel.ethCli.GetBalance(accountAddr)
		if err != nil {
			return err
		}
		if err = sel.accountRepo.UpdateBalanceById(tableName, bindInfo.ID, balance); err != nil {
			return err
		}
	}
	return nil

}
