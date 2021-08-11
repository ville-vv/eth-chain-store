package dao

import (
	"github.com/pkg/errors"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
)

// EthereumTransactionDao 以太坊交易
type EthereumTransactionDao struct {
	normalTxDao   *normalTransactionDao
	contractTxDao *contractTransactionDao
}

func NewEthereumTransactionDao(txDb DB, contractTxDb DB) *EthereumTransactionDao {
	e := &EthereumTransactionDao{
		normalTxDao: &normalTransactionDao{
			db: txDb,
			tb: newDhcpTable(txDb.GetDB(), "transaction_records"),
		},
		contractTxDao: &contractTransactionDao{
			db: contractTxDb,
			tb: newDhcpTable(contractTxDb.GetDB(), "transaction_records"),
		}}
	return e
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
	db DB
	tb *dhcpTable
}

func (sel *contractTransactionDao) createTransactionRecord(txData *model.TransactionData) error {
	tbName, err := sel.tb.TbName()
	if err != nil {
		return nil
	}
	db := sel.db.GetDB().Table(tbName)

	if err = db.Create(&model.TransactionRecord{
		BlockNumber:     txData.BlockNumber,
		BlockHash:       txData.BlockHash,
		TxHash:          txData.Hash,
		Timestamp:       txData.TimeStamp,
		ContractAddress: txData.ContractAddress,
		FromAddr:        txData.From,
		ToAddr:          txData.To,
		GasPrice:        txData.GasPrice,
		Value:           txData.Value,
		FromAddrBalance: txData.FromBalance,
		ToAddrBalance:   txData.ToBalance,
	}).Error; err != nil {
		return err
	}
	sel.tb.Inc()

	return nil
}

//==========================================================================================

type normalTransactionDao struct {
	db DB
	tb *dhcpTable
}

func (sel *normalTransactionDao) createTransactionRecord(txData *model.TransactionData) error {
	tbName, err := sel.tb.TbName()
	if err != nil {
		return nil
	}
	db := sel.db.GetDB().Table(tbName)
	if err = db.Create(&model.TransactionRecord{
		BlockNumber:     txData.BlockNumber,
		BlockHash:       txData.BlockHash,
		TxHash:          txData.Hash,
		Timestamp:       txData.TimeStamp,
		ContractAddress: txData.ContractAddress,
		FromAddr:        txData.From,
		ToAddr:          txData.To,
		GasPrice:        txData.GasPrice,
		Value:           txData.Value,
		FromAddrBalance: txData.FromBalance,
		ToAddrBalance:   "",
	}).Error; err != nil {
		return errors.Wrap(err, "create normal transaction")
	}
	sel.tb.Inc()

	return nil
}
