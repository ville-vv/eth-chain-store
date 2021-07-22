package migrate

import (
	"github.com/urfave/cli"
	"github.com/ville-vv/eth-chain-store/src/common/conf"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/vilgo/vlog"
	"github.com/ville-vv/vilgo/vstore"
	"strings"
)

func Create(_ *cli.Context) {
	businessDbCfg := conf.GetEthBusinessDbConfig()
	businessDb := vstore.MakeDBUtil(businessDbCfg)
	vlog.INFO("create db name %s", businessDbCfg.DbName)
	businessDb.CreateDB()

	ethereumDbCfg := conf.GetEthereumDbConfig()
	ethereumDb := vstore.MakeDBUtil(ethereumDbCfg)
	vlog.INFO("create db name %s", ethereumDbCfg.DbName)
	ethereumDb.CreateDB()

	contractTxDbCfg := conf.GetEthContractTransactionDbConfig()
	contractTxDb := vstore.MakeDBUtil(contractTxDbCfg)
	vlog.INFO("create db name %s", contractTxDbCfg.DbName)
	contractTxDb.CreateDB()

	transactionDbCfg := conf.GetEthTransactionDbConfig()
	transactionDb := vstore.MakeDBUtil(transactionDbCfg)
	vlog.INFO("create db name %s", transactionDbCfg.DbName)
	transactionDb.CreateDB()
}

func drop(mysqlCfg *conf.MysqlConf) {
	if strings.Contains(mysqlCfg.DbName, "prod") {
		vlog.WARN("can't drop prod db %s", mysqlCfg.DbName)
		return
	}
	utilDB := vstore.MakeDBUtil(mysqlCfg)
	vlog.INFO("drop db name %s", mysqlCfg.DbName)
	utilDB.DropDB()

}

func Drop(_ *cli.Context) {
	drop(conf.GetEthBusinessDbConfig())
	drop(conf.GetEthereumDbConfig())
	drop(conf.GetEthContractTransactionDbConfig())
	drop(conf.GetEthTransactionDbConfig())
}

func Migrate(_ *cli.Context) {
	businessDbMigrate()
	ethereumDbMigrate()
	contractTxDbDbMigrate()
	transactionDbMigrate()
}

func businessDbMigrate() {
	cfg := conf.GetEthBusinessDbConfig()
	mysqlDb := vstore.MakeDb(cfg)
	vlog.INFO("migrate db name %s", cfg.DbName)
	db := mysqlDb.GetDB()
	err := db.AutoMigrate(
		&model.Erc20ContractConfig{},
	)
	if err != nil {
		panic(err)
	}

}

func ethereumDbMigrate() {
	cfg := conf.GetEthereumDbConfig()
	mysqlDb := vstore.MakeDb(cfg)
	vlog.INFO("migrate db name %s", cfg.DbName)
	db := mysqlDb.GetDB()
	err := db.AutoMigrate(
		&model.SyncBlockConfig{},
	)
	if err != nil {
		panic(err)
	}
	db.Create(&model.SyncBlockConfig{
		KName: "SyncBlockNumber",
		Value: "0",
	})
}

func contractTxDbDbMigrate() {
	cfg := conf.GetEthContractTransactionDbConfig()
	mysqlDb := vstore.MakeDb(cfg)
	vlog.INFO("migrate db name %s", cfg.DbName)
	db := mysqlDb.GetDB()
	err := db.AutoMigrate()
	if err != nil {
		panic(err)
	}
}

func transactionDbMigrate() {
	cfg := conf.GetEthTransactionDbConfig()
	mysqlDb := vstore.MakeDb(cfg)
	vlog.INFO("migrate db name %s", cfg.DbName)
	db := mysqlDb.GetDB()
	err := db.AutoMigrate()
	if err != nil {
		panic(err)
	}
}
