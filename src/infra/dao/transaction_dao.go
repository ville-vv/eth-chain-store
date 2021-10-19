package dao

import (
	"github.com/ville-vv/eth-chain-store/src/infra/model"
)

// EthereumTransactionDao 以太坊交易
type EthereumTransactionDao struct {
	normalTxDao   *normalTransactionDao
	contractTxDao *contractTransactionDao
}

func NewEthereumTransactionDao(txDb DB, contractTxDb DB, normalTxCache *DbCache, contractTxCache *DbCache) *EthereumTransactionDao {
	e := &EthereumTransactionDao{
		normalTxDao: &normalTransactionDao{
			db:      txDb,
			tb:      newDhcpTable(txDb.GetDB(), "transaction_records"),
			dbCache: normalTxCache,
		},
		contractTxDao: &contractTransactionDao{
			db:      contractTxDb,
			tb:      newDhcpTable(contractTxDb.GetDB(), "transaction_records"),
			dbCache: contractTxCache,
		}}
	return e
}

func (sel *EthereumTransactionDao) Init() {
	if err := sel.normalTxDao.tb.Init(); err != nil {
		panic(err)
	}

	if err := sel.contractTxDao.tb.Init(); err != nil {
		panic(err)
	}
}

func (sel *EthereumTransactionDao) CreateTransactionRecord(txData *model.TransactionData) error {
	if txData.IsContractToken {
		return sel.contractTxDao.createTransactionRecord(txData)
	}
	return sel.normalTxDao.createTransactionRecord(txData)
	//return nil
}

//==========================================================================================

type contractTransactionDao struct {
	db      DB
	tb      *dhcpTable
	dbCache *DbCache
}

func (sel *contractTransactionDao) createTransactionRecord(txData *model.TransactionData) error {
	tbName, err := sel.tb.TbName()
	if err != nil {
		return nil
	}

	sel.dbCache.Insert(tbName, &model.TransactionRecord{
		BlockNumber:     txData.BlockNumber,
		BlockHash:       txData.BlockHash,
		TxHash:          txData.Hash,
		TxTime:          txData.TimeStamp,
		ContractAddress: txData.ContractAddress,
		FromAddr:        txData.From,
		ToAddr:          txData.To,
		GasPrice:        txData.GasPrice,
		Value:           txData.Value,
		FromAddrBalance: txData.FromBalance,
		ToAddrBalance:   txData.ToBalance,
	})
	// db := sel.db.GetDB().Table(tbName)
	// if err = db.Create().Error; err != nil {
	//	 return err
	// }
	sel.tb.Inc()

	return nil
}

//==========================================================================================

type normalTransactionDao struct {
	db      DB
	tb      *dhcpTable
	dbCache *DbCache
}

func (sel *normalTransactionDao) createTransactionRecord(txData *model.TransactionData) error {
	tbName, err := sel.tb.TbName()
	if err != nil {
		return nil
	}

	sel.dbCache.Insert(tbName, &model.TransactionRecord{
		BlockNumber:     txData.BlockNumber,
		BlockHash:       txData.BlockHash,
		TxHash:          txData.Hash,
		TxTime:          txData.TimeStamp,
		ContractAddress: txData.ContractAddress,
		FromAddr:        txData.From,
		ToAddr:          txData.To,
		GasPrice:        txData.GasPrice,
		Value:           txData.Value,
		FromAddrBalance: txData.FromBalance,
		ToAddrBalance:   txData.ToBalance,
	})

	//db := sel.db.GetDB().Table(tbName)
	//if err = db.Create().Error; err != nil {
	//	return errors.Wrap(err, "create normal transaction")
	//}
	sel.tb.Inc()

	return nil
}
