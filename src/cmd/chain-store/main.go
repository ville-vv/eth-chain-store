package main

import (
	"flag"
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/common/conf"
	"github.com/ville-vv/eth-chain-store/src/common/log"
	"github.com/ville-vv/eth-chain-store/src/domain/ethm"
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/dao"
	"github.com/ville-vv/eth-chain-store/src/infra/monitor"
	"github.com/ville-vv/eth-chain-store/src/infra/mqp"
	"github.com/ville-vv/eth-chain-store/src/server"
	"github.com/ville-vv/vilgo/runner"
	"github.com/ville-vv/vilgo/vlog"
	"github.com/ville-vv/vilgo/vstore"
	"os"
	"runtime"
	"time"
)

var (
	syncInterval      string
	fastSyncInterval  string
	rpcEndpoint       string
	dbUser            string
	dbPassword        string
	dbHost            string
	dbPort            string
	logFile           string
	debug             bool
	maxPullNum        int
	maxWriteNum       int
	isMaxProcs        bool
	isHelp            bool
	writeToDbInterval int
)

func cmdFlagParse() {
	flag.StringVar(&syncInterval, "si", "15", "the interval to sync latest block number")
	flag.StringVar(&fastSyncInterval, "fsi", "1000", "the interval fast to sync the block number  before  the latest ms")
	flag.StringVar(&rpcEndpoint, "rpc_endpoint", "http://localhost:8545", "eth rpc endpoint")
	flag.StringVar(&dbUser, "db_user", "", "the database user")
	flag.StringVar(&dbPassword, "db_passwd", "", "the database password")
	flag.StringVar(&dbHost, "db_host", "", "the database host")
	flag.StringVar(&dbPort, "db_port", "", "the database port")
	flag.StringVar(&logFile, "logFile", "", "the log file path and file name")
	flag.IntVar(&maxPullNum, "max_pull_num", 1, "the max thread number for sync block information from chain")
	flag.IntVar(&maxWriteNum, "max_write_num", 5, "the max thread number for write block information to db")
	flag.IntVar(&writeToDbInterval, "wi", 2, "the max thread number for write block information to db")
	flag.BoolVar(&debug, "debug", false, "open debug logs")
	flag.BoolVar(&isHelp, "help", false, "help")
	flag.BoolVar(&isMaxProcs, "max_procs", false, "the max process core")
	flag.Parse()
	fmt.Println(isHelp, debug)
	if isHelp {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if rpcEndpoint == "" {
		fmt.Println("rpc_endpoint is empty")
		flag.PrintDefaults()
		os.Exit(-1)
	}
}

func buildService() runner.Runner {
	// https://mainnet.infura.io/v3/ecc309a045134205b5c2b58481d7923d
	// https://mainnet.infura.io/v3/21628f8f9b9b423a9ea05a708016b119
	var (
		//businessDb = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthBusinessDbConfig()), "business")
		ethereumDb = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthereumDbConfig()), "ethereum")
		//contractDb    = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthContractDbConfig()), "contract")
		contractTxDb  = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthContractTransactionDbConfig()), "contract_transaction")
		transactionDb = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthTransactionDbConfig()), "transaction")
	)

	var (
		normalTxDbCache   = dao.NewDbCache(writeToDbInterval, transactionDb)
		contractTxDbCache = dao.NewDbCache(writeToDbInterval, contractTxDb)

		normalTranDao = dao.NewEthereumTransactionDao(transactionDb, contractTxDb, normalTxDbCache, contractTxDbCache)
		//erc20TokenCfgDao  = dao.NewErc20TokenConfigDao(businessDb)
		ethereumCacheDb   = dao.NewDbCache(writeToDbInterval, ethereumDb)
		ethereumDao       = dao.NewEthereumDao(ethereumDb, ethereumCacheDb)
		ethBlockNumberDao = dao.NewEthereumBlockNumberDao(ethereumDb)
		errorDao          = dao.NewSyncErrorDao(ethereumDb)
	)
	var (
		contractMng         = ethm.NewContractManager(ethm.NewEthRpcExecutor(rpcEndpoint, ""), repo.NewContractRepo(ethereumDao))
		accountMng          = ethm.NewAccountManager(ethm.NewEthRpcExecutor(rpcEndpoint, ""), repo.NewContractAccountRepo(ethereumDao), repo.NewNormalAccountRepo(ethereumDao))
		transactionWriter   = ethm.NewTransactionWriter(ethm.NewEthRpcExecutor(rpcEndpoint, ""), repo.NewTransactionRepo(normalTranDao))
		accountMngWriter    = ethm.NewRetryProcess("account", maxWriteNum, accountMng, repo.NewSyncErrorRepository(errorDao))
		contractMngWriter   = ethm.NewRetryProcess("contract", maxWriteNum, contractMng, repo.NewSyncErrorRepository(errorDao))
		transactionReWriter = ethm.NewRetryProcess("transaction", maxWriteNum, transactionWriter, repo.NewSyncErrorRepository(errorDao))

		mqPublish       = mqp.NewMDP(vlog.ERROR)
		txWriterPublish = ethm.NewEthereumPublisher(mqPublish)

		ethMng       = ethm.NewEthereumManager(ethm.NewEthRpcExecutor(rpcEndpoint, ""), txWriterPublish)
		bkNumCounter = ethm.NewSyncBlockControl(maxPullNum, ethm.NewEthRpcExecutor(rpcEndpoint, ""), repo.NewBlockNumberRepo(ethBlockNumberDao))
		serviceRun   = server.NewSyncBlockChainServiceV2(ethMng, bkNumCounter)
	)

	accountMngWriter.SetMonitor(&monitor.AccountWriteProcessNum)
	contractMngWriter.SetMonitor(&monitor.ContractWriteProcessNum)
	transactionReWriter.SetMonitor(&monitor.TxWriteProcessNum)

	mqPublish.SubScribe(accountMngWriter)
	mqPublish.SubScribe(contractMngWriter)
	mqPublish.SubScribe(transactionReWriter)

	svr := &server.Server{}
	svr.Add(serviceRun)
	svr.Add(mqPublish)
	svr.Add(accountMngWriter)
	svr.Add(contractMngWriter)
	svr.Add(transactionReWriter)
	svr.Add(ethereumCacheDb)
	svr.Add(normalTxDbCache)
	svr.Add(contractTxDbCache)

	return svr
}

func main() {
	cmdFlagParse()
	if isMaxProcs {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	log.Init() //12900839
	runner.Go(monitor.StartMonitor)
	runner.Run(buildService())
	close(conf.GlobalExitSignal)
	time.Sleep(time.Second * 5)
}
