package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/common/conf"
	"github.com/ville-vv/eth-chain-store/src/common/go_exec"
	"github.com/ville-vv/eth-chain-store/src/common/log"
	"github.com/ville-vv/eth-chain-store/src/domain/ethm"
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/dao"
	"github.com/ville-vv/eth-chain-store/src/infra/monitor"
	"github.com/ville-vv/eth-chain-store/src/server"
	"github.com/ville-vv/vilgo/runner"
	"github.com/ville-vv/vilgo/vlog"
	"github.com/ville-vv/vilgo/vstore"
	"os"
	"runtime"
	"time"
)

var (
//syncInterval      string
//fastSyncInterval  string
//rpcEndpoint       string
//dbUser            string
//dbPassword        string
//dbHost            string
//dbPort            string
//logFile           string
//debug             bool
//maxPullNum        int
//maxWriteNum       int
//isMaxProcs        bool
//isHelp            bool
//writeToDbInterval int
//txDataInHive      bool
//withTxBalance     bool
//startBlockNumber  int64 // 开始区块
//endBlockNumber    int64 // 结束区块
//saveAccount       bool
//saveContract      bool
)

func cmdFlagParse() {
	//flag.StringVar(&syncInterval, "si", "15", "the interval to sync latest block number")
	//flag.StringVar(&fastSyncInterval, "fsi", "1000", "the interval fast to sync the block number  before  the latest ms")
	flag.StringVar(&conf.RpcUrl, "rpc_url", "", "eth rpc endpoint")
	flag.StringVar(&conf.DbUser, "db_user", "", "the database user")
	flag.StringVar(&conf.DbPassword, "db_passwd", "", "the database password")
	flag.StringVar(&conf.DbHost, "db_host", "", "the database host")
	flag.StringVar(&conf.DbPort, "db_port", "", "the database port")
	flag.StringVar(&conf.LogFile, "logFile", "", "the log file path and file name")
	flag.IntVar(&conf.MaxPullNum, "max_pull_num", 1, "the max thread number for sync block information from chain")
	flag.IntVar(&conf.MaxWriteNum, "max_write_num", 5, "the max thread number for write block information to db")
	flag.IntVar(&conf.MaxBatchInsertNum, "max_insert_num", 5000, "the maximum insert number for sql data")
	flag.IntVar(&conf.WriteToDbInterval, "wi", 2, "the interval time that write to mysql from memory n/s")
	flag.BoolVar(&conf.Debug, "debug", false, "open debug logs")
	flag.BoolVar(&conf.IsHelp, "help", false, "help")
	flag.BoolVar(&conf.IsMaxProcs, "max_procs", false, "the max process core the value is true or false")
	flag.BoolVar(&conf.TxDataInHive, "txdata_inhive", false, "whether save transaction data to hive, must create the database ethereum_orc in hive db")
	flag.BoolVar(&conf.WithTxBalance, "with_tx_balance", false, "transaction data will pull balance information from eth rpc")
	flag.BoolVar(&conf.SaveAccount, "save_account", true, "save the eth account data and contract data to mysql db, ex: -save_account=false is not save")
	flag.BoolVar(&conf.SaveContract, "save_contract", true, "save the contract information to mysql db, ex: -save_contract=false is not save")
	flag.BoolVar(&conf.SaveTransaction, "save_tx", true, "save the transaction information to mysql or hive db, ex: -save_tx=false is not save")

	flag.Int64Var(&conf.StartBlockNumber, "start_number", 0, "the start block number need to sync")
	flag.Int64Var(&conf.EndBlockNumber, "end_number", 0, "the end block number need to sync ")

	flag.Parse()
	fmt.Println("Rpc:", conf.RpcUrl, "logFile:", conf.LogFile)
	if conf.IsHelp {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if conf.RpcUrl == "" {
		fmt.Println("rpc_endpoint is empty")
		flag.PrintDefaults()
		os.Exit(-1)
	}

}

func argPrint() {
	vlog.INFO("SaveAccount %t", conf.SaveAccount)
	vlog.INFO("SaveContract %t", conf.SaveContract)
	vlog.INFO("SaveTransaction %t", conf.SaveTransaction)
}

func buildService() go_exec.Runner {

	var (
		ethereumDb    = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthereumDbConfig()), "ethereum")
		contractTxDb  = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthContractTransactionDbConfig()), "contract_transaction")
		transactionDb = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthTransactionDbConfig()), "transaction")
	)

	var (
		normalTxDbCache   = dao.NewDbCache("err_data/normal_traction.sql", conf.WriteToDbInterval, transactionDb)
		contractTxDbCache = dao.NewDbCache("err_data/contract_traction.sql", conf.WriteToDbInterval, contractTxDb)
		normalTranDao     = dao.NewEthereumTransactionDao(transactionDb, contractTxDb, normalTxDbCache, contractTxDbCache)

		ethereumCacheDb   = dao.NewDbCache("err_data/ethereum_other.sql", conf.WriteToDbInterval, ethereumDb)
		ethereumDao       = dao.NewEthereumDao(ethereumDb, ethereumCacheDb)
		ethBlockNumberDao = dao.NewEthereumBlockNumberDao(ethereumDb)
		errorDao          = dao.NewSyncErrorDao(ethereumDb)
	)

	var (
		transactionRepoFactory = repo.NewTransactionRepositoryFactory(conf.TxDataInHive, conf.WriteToDbInterval, "err_data/transaction_data.sql", normalTranDao)
		transactionRepo        = transactionRepoFactory.NewTransactionRepository()
	)

	var (
		contractMng = ethm.NewContractManager(ethm.NewEthRpcExecutor(conf.RpcUrl), repo.NewContractRepo(ethereumDao))
		accountMng  = ethm.NewAccountManager(ethm.NewEthRpcExecutor(conf.RpcUrl), repo.NewContractAccountRepo(ethereumDao), repo.NewNormalAccountRepo(ethereumDao))
		//transactionWriter = ethm.NewTransactionWriter(ethm.NewEthRpcExecutor(rpcEndpoint), repo.NewTransactionRepo(normalTranDao))
		transactionWriter = ethm.NewTransactionWriter(ethm.NewEthRpcExecutor(conf.RpcUrl), transactionRepo)
		errorRepo         = repo.NewSyncErrorRepository(errorDao)

		accountMngWriter    = ethm.NewRetryProcess("account", accountMng, errorRepo)
		contractMngWriter   = ethm.NewRetryProcess("contract", contractMng, errorRepo)
		transactionReWriter = ethm.NewRetryProcess("transaction", transactionWriter, errorRepo)

		ethWriterCtl = ethm.NewEthWriterControl(conf.MaxWriteNum, accountMngWriter, contractMngWriter, transactionReWriter)

		ethMng = ethm.NewEthereumManager(ethm.NewEthRpcExecutor(conf.RpcUrl), ethWriterCtl)
		//bkNumCounter = ethm.NewSyncBlockControl(maxPullNum, ethm.NewEthRpcExecutor(rpcEndpoint), repo.NewBlockNumberRepo(ethBlockNumberDao))

		syncControl = ethm.NewSyncBlockControlWithOpt(
			&ethm.OptionConfig{
				StartBlockNumber: conf.StartBlockNumber,
				EndBlockNumber:   conf.EndBlockNumber,
				MaxSyncThreads:   conf.MaxPullNum,
				EthRpcCli:        ethm.NewEthRpcExecutor(conf.RpcUrl),
				BknRepo:          repo.NewBlockNumberRepo(ethBlockNumberDao),
			},
		)

		serviceRun = server.NewSyncBlockChainServiceV2(ethMng, syncControl)
	)

	transactionWriter.SetWithBalance(conf.WithTxBalance)

	svr := &server.Server{}
	svr.Add(serviceRun)

	// 保存账户数据
	if conf.SaveAccount {
		svr.Add(accountMngWriter)
	}

	// 保存合约数据
	if conf.SaveContract {
		svr.Add(contractMngWriter)
	}

	if conf.SaveTransaction {
		svr.Add(transactionReWriter)
		if conf.TxDataInHive {
			r, ok := transactionRepo.(runner.Runner)
			if ok {
				svr.Add(r)
			}
		} else {
			svr.Add(normalTxDbCache)
			svr.Add(contractTxDbCache)
		}
	}

	svr.Add(ethereumCacheDb)
	svr.Add(ethWriterCtl)

	return svr
}

func main() {
	cmdFlagParse()
	if conf.IsMaxProcs {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	log.Init() //12900839
	argPrint()

	go_exec.Go(monitor.StartMonitor)
	ctx, cancel := context.WithCancel(context.Background())
	go_exec.Go(func() {
		select {
		case <-conf.GlobalProgramFinishSigmal:
			cancel()
			return
		}
	})
	go_exec.Run(ctx, buildService())
	close(conf.GlobalExitSignal)
	time.Sleep(time.Second * 5)
	vlog.INFO("同步工具退出exit")
}
