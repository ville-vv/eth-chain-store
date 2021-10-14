package dao

import (
	"context"
	"fmt"
	file "github.com/ville-vv/eth-chain-store/src/common/file"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/vilgo/vlog"
	"go.uber.org/atomic"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type TransactionSQLFileDao struct {
	schema          string
	maxInsertNum    int
	dbCache         *DbCacheV2
	contractTxCnt   atomic.Int64
	ethTxCnt        atomic.Int64
	dataFile        *file.AutoSplitFile
	infoFile        *os.File
	maxDataFileSize int64 // 单位 mb
	dataFileName    string
	isExist         bool
	isFileOpen      bool
}

func NewTransactionSQLFileDao(wrInterval int, maxInsertNum int, maxDataFileSize int) *TransactionSQLFileDao {

	tx := &TransactionSQLFileDao{
		maxInsertNum:    maxInsertNum,
		maxDataFileSize: int64(maxDataFileSize),
		dataFileName:    "data/transaction_data.sql",
		schema:          "TransactionSQLFileDao",
	}
	dbCache := NewDbCacheV2(wrInterval)
	dbCache.SetExec(tx)
	tx.dbCache = dbCache

	return tx
}

func (sel *TransactionSQLFileDao) SetMaxDataFileSize(size int64) {
	sel.maxDataFileSize = size
}

func (sel *TransactionSQLFileDao) SetScheme(schema string) {
	sel.schema = schema
}

func (sel *TransactionSQLFileDao) Scheme() string {
	return sel.schema
}

func (sel *TransactionSQLFileDao) Init() error {
	// 创建自动分割文件
	dataFile, err := file.NewAutoSplitFile(sel.dataFileName, 0)
	if err != nil {
		panic(err)
	}
	dataFile.SetMaxSizeMb(sel.maxDataFileSize)
	//dataFile.FileHeaderWriteFun = func(w io.Writer) error {
	//	w.Write([]byte(fmt.Sprintf("use %s\n;", "")))
	//	return nil
	//}
	sel.dataFile = dataFile

	if err := sel.initInfoFile(); err != nil {
		vlog.ERROR("初始化数据信息错误：%v", err)
		_ = dataFile.Close()
		return err
	}
	sel.isFileOpen = true

	return nil
}

func (sel *TransactionSQLFileDao) Start() error {
	return sel.dbCache.Start()
}

func (sel *TransactionSQLFileDao) Exit(ctx context.Context) error {
	sel.isExist = false
	sel.dbCache.Exit(ctx)
	if sel.isFileOpen {
		sel.saveInfo()
		sel.infoFile.Close()
		sel.dataFile.Close()
	}
	return nil
}

func (sel *TransactionSQLFileDao) Exec(tableName string, record []interface{}) error {
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
		_, err := sel.dataFile.WriteString(insertSql + ";\n")
		if err != nil {
			vlog.ERROR("save data to hive db table %s len:%d error %s", tableName, len(spList), err.Error())
			return err
		}
	}
	return nil
}

func (sel *TransactionSQLFileDao) initInfoData(fileName string) error {
	// 读取原来保存的信息
	ff, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer ff.Close()

	data, err := ioutil.ReadAll(ff)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		sel.contractTxCnt.Store(0)
		sel.ethTxCnt.Store(0)
		return nil
	}

	vlog.INFO("读取到信息：%s", string(data))

	var ethTxCntStr, contractTxCntStr string
	strList := strings.Split(string(data), ",")
	if len(strList) > 1 {
		ethTxCntStr = strList[0]
		contractTxCntStr = strList[1]
	}
	ethTxCnt, _ := strconv.ParseInt(ethTxCntStr, 10, 64)
	vlog.INFO("eth transaction start id %d", ethTxCnt)
	sel.ethTxCnt.Store(ethTxCnt)
	contractTxCnt, _ := strconv.ParseInt(contractTxCntStr, 10, 64)
	vlog.INFO("contract transaction start id %d", ethTxCnt)
	sel.contractTxCnt.Store(contractTxCnt)
	return nil
}

func (sel *TransactionSQLFileDao) initInfoFile() error {
	dirPath := path.Dir(sel.dataFileName)
	if dirPath == "" {
		dirPath = "./"
	}
	fileName := path.Join(dirPath, "data_info")

	err := sel.initInfoData(fileName)
	if err != nil {
		return err
	}

	ff, err := os.OpenFile(fileName, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	sel.infoFile = ff

	go func() {
		for {
			time.Sleep(time.Second)
			if sel.isExist {
				return
			}
			sel.saveInfo()
		}
	}()
	return nil
}

func (sel *TransactionSQLFileDao) saveInfo() {
	//sel.infoFile.Truncate(0)
	_, _ = sel.infoFile.Seek(0, io.SeekStart)
	_, _ = sel.infoFile.WriteString(fmt.Sprintf("%d,%d", sel.ethTxCnt.Load(), sel.contractTxCnt.Load()))
}

func (sel *TransactionSQLFileDao) CreateTransactionRecord(txData *model.TransactionData) error {
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
