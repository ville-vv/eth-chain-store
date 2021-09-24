package migrate

import (
	"github.com/urfave/cli"
	"github.com/ville-vv/eth-chain-store/src/common/conf"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/vilgo/vlog"
	"github.com/ville-vv/vilgo/vstore"
	"gorm.io/gorm"
	"strings"
)

func Create(_ *cli.Context) {
	//businessDbCfg := conf.GetEthBusinessDbConfig()
	//businessDb := vstore.MakeDBUtil(businessDbCfg)
	//vlog.INFO("create db name %s", businessDbCfg.DbName)
	//businessDb.CreateDB()

	create := func(dbConf *conf.MysqlConf) {
		db := vstore.MakeDBUtil(dbConf)
		vlog.INFO("create db name %s", dbConf.DbName)
		db.CreateDB()
	}

	create(conf.GetEthereumDbConfig())
	create(conf.GetEthContractTransactionDbConfig())
	create(conf.GetEthTransactionDbConfig())
	create(conf.GetEthereumHiveMapDbConfig())
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
	// drop(conf.GetEthBusinessDbConfig())
	drop(conf.GetEthereumDbConfig())
	drop(conf.GetEthContractTransactionDbConfig())
	drop(conf.GetEthTransactionDbConfig())
	drop(conf.GetEthereumHiveMapDbConfig())
}

func Migrate(_ *cli.Context) {
	//businessDbMigrate()
	ethereumDbMigrate()
	contractTxDbDbMigrate()
	transactionDbMigrate()
	ethereumHiveMapDbMigrate()

}

func ethereumDbMigrate() {
	cfg := conf.GetEthereumDbConfig()
	mysqlDb := vstore.MakeDb(cfg)
	vlog.INFO("migrate db name %s", cfg.DbName)
	db := mysqlDb.GetDB().Debug()
	err := db.AutoMigrate(
		&model.SyncBlockConfig{},
		&model.SplitTableInfo{},
		&model.ContractAccountBind{},
		&model.EthereumAccount{},
		//&model.ContractAddressRecord{},
		&model.SyncErrorRecord{},
		&model.ContractAddressRecord{},
	)
	if err != nil {
		panic(err)
	}
}

func contractTxDbDbMigrate() {
	cfg := conf.GetEthContractTransactionDbConfig()
	mysqlDb := vstore.MakeDb(cfg)
	vlog.INFO("migrate db name %s", cfg.DbName)
	db := mysqlDb.GetDB().Debug()
	err := db.AutoMigrate(
		&model.SplitTableInfo{},
		&model.TransactionRecord{},
	)
	if err != nil {
		panic(err)
	}
	//createMyISAMTable(db, &model.TransactionRecord{})

}

func transactionDbMigrate() {
	cfg := conf.GetEthTransactionDbConfig()
	mysqlDb := vstore.MakeDb(cfg)
	vlog.INFO("migrate db name %s", cfg.DbName)
	db := mysqlDb.GetDB().Debug()
	err := db.AutoMigrate(
		&model.SplitTableInfo{},
		&model.TransactionRecord{},
	)

	//db.Set("gorm:table_options", "ENGINE=MyISAM").AutoMigrate(&model.TransactionRecord{})
	//db.Set("gorm:table_options", "ENGINE=MyISAM").AutoMigrate(&model.TransactionRecord{})

	if err != nil {
		panic(err)
	}
}

func ethereumHiveMapDbMigrate() {
	cfg := conf.GetEthereumHiveMapDbConfig()
	mysqlDb := vstore.MakeDb(cfg)
	vlog.INFO("migrate db name %s", cfg.DbName)
	db := mysqlDb.GetDB().Debug()
	err := db.AutoMigrate(
		&model.SyncBlockConfig{},
		&model.SplitTableInfo{},
		&model.ContractAddressRecord{},
		&model.ContractAccountBind{},
		&model.EthereumAccount{},
		&model.SyncErrorRecord{},
	)
	if err != nil {
		panic(err)
	}
}

func createMyISAMTable(db *gorm.DB, tb ...interface{}) {
	err := db.Set("gorm:table_options", "ENGINE=MyISAM;").AutoMigrate(tb...)
	if err != nil {
		panic(err)
	}
}
