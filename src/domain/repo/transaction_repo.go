package repo

import (
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/common/conf"
	"github.com/ville-vv/eth-chain-store/src/infra/dao"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
)

type TransactionRepository interface {
	CreateTransactionRecord(txData *model.TransactionData) error
}

type TransactionRepositoryFactory struct {
	SaveType          string
	WriteToDbInterval int
	ErrFile           string
	EthereumTranDao   *dao.EthereumTransactionDao
}

func NewTransactionRepositoryFactory(inHive bool, writeToDbInterval int, errFile string, ethereumTranDao *dao.EthereumTransactionDao) *TransactionRepositoryFactory {
	saveType := ""
	if inHive {
		saveType = "InHive"
	}
	return &TransactionRepositoryFactory{SaveType: saveType, WriteToDbInterval: writeToDbInterval, ErrFile: errFile, EthereumTranDao: ethereumTranDao}
}

func (sel *TransactionRepositoryFactory) NewTransactionRepository() TransactionRepository {
	switch sel.SaveType {
	case "InHive":
		fmt.Println("inHive")
		return dao.NewTransactionHiveDao(sel.ErrFile, conf.GetHiveEthereumDb(), sel.WriteToDbInterval, conf.MaxBatchInsertNum)
	default:
		return NewTransactionRepo(sel.EthereumTranDao)
	}
}

// 交易数据处理
type TransactionRepo struct {
	normalTranDao *dao.EthereumTransactionDao
}

func NewTransactionRepo(normalTranDao *dao.EthereumTransactionDao) *TransactionRepo {
	return &TransactionRepo{normalTranDao: normalTranDao}
}

func (sel *TransactionRepo) CreateTransactionRecord(txData *model.TransactionData) error {
	return sel.normalTranDao.CreateTransactionRecord(txData)
}
