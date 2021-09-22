package dao

import (
	"github.com/ville-vv/eth-chain-store/src/common/hive"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/vilgo/vlog"
	"github.com/ville-vv/vilgo/vstore"
	"os"
)

type EthereumMapHive struct {
	db      vstore.DB
	hiveCli *hive.HiveCLI
	dbCache *HiveDbCache
	errFile *os.File
}

func (sel *EthereumMapHive) QueryContractInfo(addr string, contractInfo *model.ContractAddressRecord) error {
	return sel.db.GetDB().Where("address=?").First(contractInfo).Error
}

func (sel *EthereumMapHive) CreateContractRecord(contractCtx *model.ContractContent) error {
	return sel.dbCache.Insert("contract_address_records", &model.ContractAddressRecord{ContractContent: *contractCtx})
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
