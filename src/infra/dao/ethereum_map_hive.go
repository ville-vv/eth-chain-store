package dao

import (
	"context"
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/common/go_exec"
	"github.com/ville-vv/eth-chain-store/src/common/hive"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/vilgo/vfile"
	"github.com/ville-vv/vilgo/vlog"
	"github.com/ville-vv/vilgo/vstore"
	"gorm.io/gorm"
	"os"
	"path"
)

const (
	TbTransactionRecords         = "transaction_records_orc"
	TbContractTransactionRecords = "contract_transaction_records_orc"
)

type EthereumMapHive struct {
	db      vstore.DB
	hiveCli *hive.HiveCLI
	dbCache *DbCacheV2
	errFile *os.File
}

func NewEthereumMapHive(errFile string, db vstore.DB, hiveCli *hive.HiveCLI, wrInterval int) *EthereumMapHive {
	if wrInterval <= 0 {
		wrInterval = 1
	}
	var err error
	e := &EthereumMapHive{db: db, hiveCli: hiveCli}
	e.dbCache = NewDbCacheV2WithMaxCache(wrInterval, 2)
	e.dbCache.SetExec(e)
	dirPath := path.Dir(errFile)
	if !vfile.PathExists(dirPath) {
		err := os.Mkdir(dirPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	e.errFile, err = os.OpenFile(errFile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	return e
}

func (sel *EthereumMapHive) Start() error {
	go_exec.Go(func() {
		sel.dbCache.Start()
	})
	return nil
}

func (sel *EthereumMapHive) Exit(ctx context.Context) error {
	return sel.dbCache.Exit(context.Background())
}

func (sel *EthereumMapHive) QueryContractInfo(addr string, contractInfo *model.ContractAddressRecord) error {
	return sel.db.GetDB().Where("address=?", addr).First(contractInfo).Error
}

func (sel *EthereumMapHive) CreateContractRecord(contractCtx *model.ContractContent) error {
	return sel.dbCache.InsertAndWait("contract_address_records", &model.ContractAddressRecord{ContractContent: *contractCtx})
}

func (sel *EthereumMapHive) Exec(tableName string, record []interface{}) error {
	insertSql := BatchInsertToSqlStr(tableName, record)
	db := sel.db.GetDB().Begin()
	defer db.Rollback()
	err := sel.db.GetDB().Exec(insertSql).Error
	if err != nil {
		vlog.ERROR("save data to data table %s len:%d error %s", tableName, len(record), err.Error())
		_, _ = sel.errFile.WriteString(insertSql + ";\n")
	}
	return db.Commit().Error
}

func (sel *EthereumMapHive) GetTxRecordAroundBlockNo(tbName string, blockNo int64, blockSize int64) ([]*model.TransactionRecord, error) {
	list := make([]*model.TransactionRecord, 0, 10000)
	stm := fmt.Sprintf("select block_number,block_hash,tx_hash,tx_time,contract_address,from_addr, distinct to_addr from "+tbName+
		"where block_number between '%d' and '%d'",
		blockNo,
		blockNo+blockSize)
	stm = "select distinct to_addr ,block_number,block_hash,tx_hash,tx_time,contract_address,from_addr from " + tbName + " limit 1000"
	stm = "select to_addr ,block_number,block_hash,tx_hash,tx_time,contract_address,from_addr from " + tbName + " limit 1000"
	fmt.Println(stm)
	err := sel.hiveCli.Find(stm, &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (sel *EthereumMapHive) UpdateFinishInfo(typ string, content string) error {
	db := sel.db.GetDB()
	return db.Model(&model.SyncBlockConfig{}).Where("k_name=?", typ).Update("value", content).Error
}

func (sel *EthereumMapHive) QueryCursorInfo(typ string) (string, error) {
	db := sel.db.GetDB().Debug()
	var sbc = model.SyncBlockConfig{}
	err := db.Model(&sbc).Where("k_name=?", typ).First(&sbc).Error
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			return "", err
		}
		return "", nil
	}
	return sbc.Value, nil
}

func (sel *EthereumMapHive) CreateCursorInfo(typ string, content string) error {
	db := sel.db.GetDB().Debug()
	return db.Create(&model.SyncBlockConfig{KName: typ, Value: content}).Error
}
