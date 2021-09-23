package server

import (
	"context"
	"github.com/ville-vv/eth-chain-store/src/common/conf"
	"github.com/ville-vv/eth-chain-store/src/common/go_exec"
	"github.com/ville-vv/eth-chain-store/src/common/hive"
	"github.com/ville-vv/eth-chain-store/src/common/log"
	"github.com/ville-vv/eth-chain-store/src/domain/async"
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/dao"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/vilgo/vstore"
	"os"
	"testing"
)

func buildService() go_exec.Runner {
	os.Setenv("HIVE2_USER_NAME", "hadoop")
	os.Setenv("HIVE2_PASSWORD", "")
	os.Setenv("HIVE2_HOST", "172.16.16.155")
	os.Setenv("HIVE2_PORT", "10000")
	os.Setenv("HIVE2_AUTH_MODEL", "NOSASL")

	var (
		rpcEndpoint = "http://172.16.16.115:8545"
	)
	dbMysql := vstore.MakeDb(conf.GetEthereumHiveMapDbConfig())
	hiveConfig := conf.GetHiveEthereumDb()
	hiveConfig.DbName = "etherum_orc"
	hiveCli, err := hive.New(hiveConfig)
	if err != nil {
		panic(err)
	}

	var (
		ethereumCacheDb = dao.NewDbCache("err_data/ethereum_map_hive_01.sql", 1, dbMysql)
		ethereumDao     = dao.NewEthereumDao(dbMysql, ethereumCacheDb)

		cursorRepo           = dao.NewEthereumMapHive("err_data/ethereum_map_hive_02.sql", dbMysql, hiveCli, 1)
		errorDao             = dao.NewSyncErrorDao(dbMysql)
		errorRepo            = repo.NewSyncErrorRepository(errorDao)
		latestBlockNumGetter = async.NewLatestBlockNumberCache(ethrpc.NewClient(rpcEndpoint))
	)

	var (
		contractDataCursor = async.NewDataCursorAggregate(async.CursorTypeContractTx, cursorRepo)
		contractAccount    = async.NewContractAccountService(ethrpc.NewClient(rpcEndpoint), repo.NewContractAccountRepo(ethereumDao))
		contractPcr        = async.NewDataProcessorCtl(contractAccount, contractDataCursor, errorRepo, ethrpc.NewClient(rpcEndpoint))
	)

	var (
		ethDataCursor = async.NewDataCursorAggregate(async.CursorTypeEthereumTx, cursorRepo)
		ethAccount    = async.NewEthAccountService(ethrpc.NewClient(rpcEndpoint), repo.NewNormalAccountRepo(ethereumDao))
		ethAccountPcr = async.NewDataProcessorCtl(ethAccount, ethDataCursor, errorRepo, ethrpc.NewClient(rpcEndpoint))
	)
	timerSvr := NewTimerServer()
	timerSvr.Add(latestBlockNumGetter, cursorRepo, contractPcr, ethAccountPcr)

	return timerSvr
}

func TestNewTimerServer(t *testing.T) {
	log.Init()
	start := buildService()
	start.Start(context.Background())
	select {}
}
