package repo

import (
	"github.com/ville-vv/eth-chain-store/src/infra/dao"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
)

type TransactionRepository interface {
	CreateTransactionRecord(txData *model.TransactionData) error
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
