package dao

import (
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/common/conf"
	"github.com/ville-vv/eth-chain-store/src/common/hive"
	"github.com/ville-vv/vilgo/vstore"
	"os"
	"testing"
)

func TestEthereumMapHive_GetTxRecordAroundBlockNo(t *testing.T) {
	os.Setenv("HIVE2_USER_NAME", "hadoop")
	os.Setenv("HIVE2_PASSWORD", "")
	os.Setenv("HIVE2_HOST", "172.16.16.155")
	os.Setenv("HIVE2_PORT", "10000")
	os.Setenv("HIVE2_AUTH_MODEL", "NOSASL")
	hiveConfig := conf.GetHiveEthereumDb()
	hiveConfig.DbName = "etherum_orc"

	dbMysql := vstore.MakeDb(conf.GetEthereumHiveMapDbConfig())
	hiveCli, err := hive.New(hiveConfig)
	if err != nil {
		t.Error(err)
		return
	}
	cursorRepo := NewEthereumMapHive("err_data/ethereum_map_hive_02.sql", dbMysql, hiveCli, 1)
	dataList, err := cursorRepo.GetTxRecordAroundBlockNo("transaction_records_orc", 0, 5)
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < len(dataList); i++ {
		fmt.Println(dataList[i].TxHash, dataList[i].ToAddr, dataList[i].CreatedAt)
	}

}
