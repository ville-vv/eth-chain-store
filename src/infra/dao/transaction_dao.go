package dao

import (
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"gorm.io/gorm"
	"sync"
)

type dhcpTable struct {
	cntTable string
	id       int64
	db       *gorm.DB
	lock     sync.Mutex
	maxNum   int64
	counter  int64 // 计数器
}

func newDhcpTable(db *gorm.DB) *dhcpTable {
	tb := &dhcpTable{
		db:      db,
		lock:    sync.Mutex{},
		maxNum:  100000000,
		counter: 0,
	}
	err := tb.intCntTxTable()
	if err != nil {
		panic(err)
	}
	return tb
}

// intCntTxTable 初始化
func (sel *dhcpTable) intCntTxTable() error {
	db := sel.db
	tb := &model.TransactionAllTable{}
	err := db.Model(tb).Select("id, table_name").Last(tb).Error
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			return nil
		}
		err = sel.createTxTable(sel.id)
		if err != nil {
			return nil
		}
		sel.id = 1
	} else {
		sel.cntTable = tb.TableName
		sel.id = tb.ID
	}
	return nil
}

func (sel *dhcpTable) createTxTable(id int64) error {
	db := sel.db.Begin()
	db.Rollback()
	tbName := fmt.Sprintf("transaction_record_%d", id)
	_, err := db.Exec(fmt.Sprintf("create table %s like transaction_table_templates", tbName)).Rows()
	if err != nil {
		return err
	}
	tb := &model.TransactionAllTable{TableName: tbName}
	err = db.Create(tb).Error
	if err != nil {
		return err
	}
	sel.id = tb.ID
	sel.cntTable = tb.TableName
	return err
}

// 表格覆盖
func (sel *dhcpTable) RecoverTbName() (string, error) {
	sel.lock.Lock()
	if sel.counter > sel.maxNum {
		err := sel.createTxTable(sel.id)
		if err != nil {
			return "", err
		}
		sel.id++
		sel.counter = 0
	}
	sel.lock.Unlock()
	return sel.cntTable, nil
}

// EthereumTransactionDao 以太坊交易
type EthereumTransactionDao struct {
	normalTxDao   *normalTransactionDao
	contractTxDao *contractTransactionDao
}

func NewEthereumTransactionDao(txDb DB, contractTxDb DB) *EthereumTransactionDao {
	e := &EthereumTransactionDao{
		normalTxDao: &normalTransactionDao{
			db: txDb,
			tb: newDhcpTable(txDb.GetDB()),
		},
		contractTxDao: &contractTransactionDao{
			db: contractTxDb,
			tb: newDhcpTable(contractTxDb.GetDB()),
		}}
	return e
}

func (sel *EthereumTransactionDao) CreateTransactionRecord(txData *model.TransactionData) error {
	if txData.IsContractToken {
		return sel.contractTxDao.createTransactionRecord(txData)
	}
	fmt.Println("写入交易数据", txData)
	return nil
}

//==========================================================================================

type contractTransactionDao struct {
	db DB
	tb *dhcpTable
}

func (sel *contractTransactionDao) createTransactionRecord(txData *model.TransactionData) error {
	tbName, err := sel.tb.RecoverTbName()
	if err != nil {
		return nil
	}
	db := sel.db.GetDB().Table(tbName)
	return db.Create(&model.TransactionTableTemplate{
		TxHash: txData.Hash,
		//Address:  txData.ContractAddress,
		FromAddr: txData.From,
		ToAddr:   txData.To,
		Gas:      txData.GasPrice,
		Value:    txData.Value,
		Balance:  txData.Balance,
	}).Error
}

//==========================================================================================

type normalTransactionDao struct {
	db DB
	tb *dhcpTable
}

func (sel *normalTransactionDao) createTransactionRecord(txData *model.TransactionData) error {
	tbName, err := sel.tb.RecoverTbName()
	if err != nil {
		return nil
	}
	db := sel.db.GetDB().Table(tbName)
	return db.Create(txData).Error
}
