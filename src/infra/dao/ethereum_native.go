package dao

import (
	"github.com/ville-vv/eth-chain-store/src/common/conf"
	"github.com/ville-vv/vilgo/vstore"
)

type EthereumNativeDao struct {
	db DB
}

func NewEthereumNativeDao() *EthereumNativeDao {
	dbCfg := conf.NewMysqlConf()
	dbCfg.DbName = "ethereum_relative"
	return &EthereumNativeDao{db: NewMysqlDB(vstore.MakeDb(dbCfg), dbCfg.DbName)}
}

func (sel *EthereumNativeDao) Get() {
}
