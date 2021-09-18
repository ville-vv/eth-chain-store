package async

import (
	"github.com/ville-vv/eth-chain-store/src/common/utils"
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"time"
)

// EthAccountService
type EthAccountService struct {
	rpcCli      ethrpc.EthRPC
	accountRepo repo.EthAccountRepository
}

func (sel *EthAccountService) Process(record *model.TransactionRecord) error {
	addr := record.ToAddr
	ok, err := sel.accountRepo.IsAccountExist(addr)
	if err != nil {
		return err
	}
	if ok {
		timeStmp, err := utils.ParseLocal(record.TxTime)
		if err != nil {
			return nil
		}
		if !timeStmp.After(time.Now().Add(-1 * time.Hour)) {
			return nil
		}
		//存在，判断是否需要更新
		balance, err := sel.rpcCli.GetBalance(addr)
		if err != nil {
			return err
		}
		err = sel.accountRepo.UpdateAccountByAddr(addr, map[string]interface{}{
			"balance": balance,
		})
		return err
	}

	balance, err := sel.rpcCli.GetBalance(addr)
	if err != nil {
		return err
	}

	// 不存在就创建
	return sel.accountRepo.CreateAccount(&model.EthereumAccount{
		Address:     addr,
		FirstTxTime: record.TxTime,
		FirstTxHash: record.TxHash,
		Balance:     balance,
	})
}
