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
	hiveCliContractAccount, err := hive.New(hiveConfig)
	if err != nil {
		panic(err)
	}

	hiveCliEthereumAccount, err := hive.New(hiveConfig)
	if err != nil {
		panic(err)
	}

	var (
		ethereumCacheDb = dao.NewDbCache("err_data/ethereum_map_hive_01.sql", 1, dbMysql)
		ethereumDao     = dao.NewEthereumDao(dbMysql, ethereumCacheDb)

		errorDao             = dao.NewSyncErrorDao(dbMysql)
		errorRepo            = repo.NewSyncErrorRepository(errorDao)
		latestBlockNumGetter = async.NewLatestBlockNumberCache(ethrpc.NewClient(rpcEndpoint))
	)
	var (
		// ethereum  Process Control
		ethcursorRepo = dao.NewEthereumMapHive("err_data/ethereum_map_hive_03.sql", dbMysql, hiveCliEthereumAccount, 1)
		ethDataCursor = async.NewDataCursorAggregate(async.CursorTypeEthereumTx, ethcursorRepo)
		ethAccountPcr = async.NewDataProcessorCtl(ethDataCursor, errorRepo, ethrpc.NewClient(rpcEndpoint))

		// contract Process Control
		contractcursorRepo = dao.NewEthereumMapHive("err_data/ethereum_map_hive_02.sql", dbMysql, hiveCliContractAccount, 1)
		contractDataCursor = async.NewDataCursorAggregate(async.CursorTypeContractTx, contractcursorRepo)
		contractPcr        = async.NewDataProcessorCtl(contractDataCursor, errorRepo, ethrpc.NewClient(rpcEndpoint))
	)
	contractDataCursor.Init()
	ethDataCursor.Init()

	var (
		ethAccount      = async.NewEthAccountService(ethrpc.NewClient(rpcEndpoint), repo.NewNormalAccountRepo(ethereumDao))
		contractAccount = async.NewContractAccountService(ethrpc.NewClient(rpcEndpoint), repo.NewContractAccountRepo(ethereumDao))
		contractInfo    = async.NewContractService(ethrpc.NewClient(rpcEndpoint), repo.NewContractRepo(ethereumDao))
	)
	contractPcr.AddProcess(contractAccount)
	contractPcr.AddProcess(contractInfo)
	ethAccountPcr.AddProcess(ethAccount)

	timerSvr := NewTimerServer()
	timerSvr.Add(ethcursorRepo, contractcursorRepo)
	timerSvr.Add(latestBlockNumGetter, ethereumCacheDb, contractPcr, ethAccountPcr)

	return timerSvr
}

func TestNewTimerServer(t *testing.T) {
	log.Init()
	start := buildService()

	go_exec.Run(context.Background(), start)

	select {}
}
