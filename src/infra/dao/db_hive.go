package dao

import (
	"context"
	"github.com/ville-vv/eth-chain-store/src/common/hive"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/vilgo/vfile"
	"github.com/ville-vv/vilgo/vlog"
	"go.uber.org/atomic"
	"os"
	"path"
)

type TransactionHiveDao struct {
	dbCache *HiveDbCache
	hiveCli *hive.HiveCLI
	errFile *os.File

	contractTxCnt atomic.Int64
	ethTxCnt      atomic.Int64
}

func (sel *TransactionHiveDao) Scheme() string {
	return "TransactionHiveDao"
}

func (sel *TransactionHiveDao) Init() error {
	vlog.INFO("正在统计 transaction_records 总数")
	sel.ethTxCnt.Store(sel.hiveCli.Count("transaction_records"))
	vlog.INFO("正在统计 contract_transaction_records 总数")
	sel.contractTxCnt.Store(sel.hiveCli.Count("contract_transaction_records"))
	return nil
}

func (sel *TransactionHiveDao) Start() error {
	return sel.dbCache.Start()
}

func (sel *TransactionHiveDao) Exit(ctx context.Context) error {
	sel.dbCache.Exit(ctx)
	sel.errFile.Close()
	return nil
}

func NewTransactionHiveDao(errFile string, option hive.HiveConfigOption) *TransactionHiveDao {
	var err error
	thd := &TransactionHiveDao{}
	thd.dbCache = NewHiveDbCache(thd)
	thd.hiveCli, err = hive.New(option)
	if err != nil {
		panic(err)
	}

	dirPath := path.Dir(errFile)
	if !vfile.PathExists(dirPath) {
		err := os.Mkdir(dirPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	thd.errFile, err = os.OpenFile(errFile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	return thd
}

func (sel *TransactionHiveDao) Exec(tableName string, record []interface{}) error {
	insertSql := BatchInsertToSqlStr(tableName, record)
	//fmt.Println(insertSql)
	err := sel.hiveCli.Exec(insertSql)
	if err != nil {
		vlog.ERROR("save data to hive db table %s len:%d error %s", tableName, len(record), err.Error())
		_, _ = sel.errFile.WriteString(insertSql + ";\n")
	}
	return nil
}

func (sel *TransactionHiveDao) CreateTransactionRecord(txData *model.TransactionData) error {
	var id = int64(1)
	var tbName = "transaction_records"
	if txData.IsContractToken {
		tbName = "contract_transaction_records"
		sel.contractTxCnt.Inc()
		id = sel.contractTxCnt.Load()
	} else {
		sel.ethTxCnt.Inc()
		id = sel.ethTxCnt.Load()
	}

	return sel.dbCache.Insert(tbName, &model.TransactionRecord{
		ID:              id,
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
}
