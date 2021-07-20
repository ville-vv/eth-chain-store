package repo

import (
	"github.com/ville-vv/eth-chain-store/src/infra/dao"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
)

// 交易数据处理
type TransactionRepo struct {
	normalTranDao *dao.EthereumTransactionDao
}

func (sel *TransactionRepo) CreateTransactionRecord(txData *model.TransactionData) error {
	return sel.normalTranDao.CreateTransactionRecord(txData)
}
