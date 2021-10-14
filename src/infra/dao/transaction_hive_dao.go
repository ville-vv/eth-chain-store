package dao

import (
	"context"
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/common/hive"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/vilgo/vfile"
	"github.com/ville-vv/vilgo/vlog"
	"go.uber.org/atomic"
	"os"
	"path"
	"time"
)

type TransactionHiveDao struct {
	dbCache *DbCacheV2
	hiveCli *hive.HiveCLI
	errFile *os.File

	maxInsertNum int

	contractTxCnt atomic.Int64
	ethTxCnt      atomic.Int64
}

func (sel *TransactionHiveDao) Scheme() string {
	return "TransactionHiveDao"
}

func (sel *TransactionHiveDao) Init() error {
	vlog.INFO("正在统计 hive transaction_records 总数")
	countor := sel.hiveCli.Count("transaction_records")
	vlog.INFO("总数为：%d", countor)
	sel.ethTxCnt.Store(countor)
	vlog.INFO("正在统计 hive contract_transaction_records 总数")
	countor = sel.hiveCli.Count("contract_transaction_records")
	vlog.INFO("总数为：%d", countor)
	sel.contractTxCnt.Store(countor)
	return nil
}

func (sel *TransactionHiveDao) Start() error {
	vlog.INFO("Start transaction data save to hive thread ")
	return sel.dbCache.Start()
}

func (sel *TransactionHiveDao) Exit(ctx context.Context) error {
	sel.dbCache.Exit(ctx)
	vlog.INFO("hive transaction data thread  exit")
	time.Sleep(time.Second * 1)
	sel.errFile.Close()
	sel.hiveCli.Close()
	return nil
}

func NewTransactionHiveDao(errFile string, option hive.HiveConfigOption, wrInterval int, maxInsertNum int) *TransactionHiveDao {
	if maxInsertNum == 0 {
		maxInsertNum = 500000
	}
	var err error
	thd := &TransactionHiveDao{}
	thd.dbCache = NewDbCacheV2WithMaxCache(wrInterval, 1000)
	thd.dbCache.SetExec(thd)
	thd.hiveCli, err = hive.New(option)
	thd.maxInsertNum = maxInsertNum
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
	_, _ = thd.errFile.WriteString(fmt.Sprintf("use %s;\n", option.GetDBName()))

	return thd
}

func (sel *TransactionHiveDao) Exec(tableName string, record []interface{}) error {
	spSize := sel.maxInsertNum
	num := len(record)
	sp := (num / spSize) + 1
	var spList []interface{}
	for i := 0; i < sp; i++ {

		if i+1 == sp {
			spList = record[i*spSize:]
		} else {
			spList = record[i*spSize : (i+1)*spSize]
		}

		insertSql := BatchInsertToSqlNoTitle(tableName, spList)
		err := sel.hiveCli.Exec(insertSql)
		if err != nil {
			vlog.ERROR("save data to hive db table %s len:%d error %s", tableName, len(spList), err.Error())
			_, _ = sel.errFile.WriteString(insertSql + ";\n")
		}
	}

	//insertSql := BatchInsertToSqlStrNeedID(tableName, record)
	//err := sel.hiveCli.Exec(insertSql)
	//if err != nil {
	//	vlog.ERROR("save data to hive db table %s len:%d error %s", tableName, len(record), err.Error())
	//	_, _ = sel.errFile.WriteString(insertSql + ";\n")
	//}
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

	return sel.dbCache.InsertAndWait(tbName, &model.TransactionRecord{
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
