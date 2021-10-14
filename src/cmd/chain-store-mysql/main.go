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
	flag.IntVar(&conf.MaxBatchInsertNum, "max_insert_num", 4000, "the maximum insert number for sql data")
	flag.IntVar(&conf.MaxSqlFileSize, "max_sql_file_size", 1000, "save type is InSqlFile, using to limit the single sql data file size unit mb")
	flag.IntVar(&conf.WriteToDbInterval, "wi", 2, "the interval time that write to mysql from memory n/s")
	flag.BoolVar(&conf.Debug, "debug", false, "open debug logs")
	flag.BoolVar(&conf.IsHelp, "help", false, "help")
	flag.BoolVar(&conf.IsMaxProcs, "max_procs", false, "the max process core the value is true or false")
	flag.BoolVar(&conf.TxDataInHive, "txdata_inhive", false, "whether save transaction data to hive, must create the database ethereum_orc in hive db")
	flag.BoolVar(&conf.WithTxBalance, "with_tx_balance", false, "transaction data will pull balance information from eth rpc")
	flag.BoolVar(&conf.SaveAccount, "save_account", true, "save the eth account data and contract data to mysql db, ex: -save_account=false is not save")
	flag.BoolVar(&conf.SaveContract, "save_contract", true, "save the contract information to mysql db, ex: -save_contract=false is not save")
	flag.BoolVar(&conf.SaveTransaction, "save_tx", true, "save the transaction information to mysql or hive db, ex: -save_tx=false is not save")
	flag.StringVar(&conf.SaveType, "save_type", "InMysql", "where are you want to save the transaction, default is save to mysql, the option value is InHive, InSqlFile")

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
	vlog.INFO("SaveAccount: %t", conf.SaveAccount)
	vlog.INFO("SaveContract: %t", conf.SaveContract)
	vlog.INFO("SaveTransaction: %t", conf.SaveTransaction)
	vlog.INFO("SaveType: %s", conf.SaveType)
}

func buildService() go_exec.Runner {

	var (
		ethereumDb        = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthereumDbConfig()), "ethereum")
		ethereumCacheDb   = dao.NewDbCache("err_data/ethereum_other.sql", conf.WriteToDbInterval, ethereumDb)
		ethereumDao       = dao.NewEthereumDao(ethereumDb, ethereumCacheDb)
		ethBlockNumberDao = dao.NewEthereumBlockNumberDao(ethereumDb)
		errorDao          = dao.NewSyncErrorDao(ethereumDb)
		errorRepo         = repo.NewSyncErrorRepository(errorDao)
	)

	var (
		ethWriterCtl = ethm.NewEthWriterControl(conf.MaxWriteNum)
		ethMng       = ethm.NewEthereumManager(ethm.NewEthRpcExecutor(conf.RpcUrl), ethWriterCtl)
		syncControl  = ethm.NewSyncBlockControlWithOpt(
			&ethm.OptionConfig{
				StartBlockNumber: conf.StartBlockNumber,
				EndBlockNumber:   conf.EndBlockNumber,
				MaxSyncThreads:   conf.MaxPullNum,
				EthRpcCli:        ethm.NewEthRpcExecutor(conf.RpcUrl),
				BknRepo:          repo.NewBlockNumberRepo(ethBlockNumberDao),
			},
		)
		serviceRun = server.NewSyncBlockChainServiceV2(ethMng, syncControl)
		svr        = &server.Server{}
	)

	svr.Add(serviceRun)
	svr.Add(ethWriterCtl)
	svr.Add(ethereumCacheDb)

	if conf.SaveContract {
		ethereumDao.InitContractRecordTb()
		accountMng := ethm.NewAccountManager(ethm.NewEthRpcExecutor(conf.RpcUrl), repo.NewContractAccountRepo(ethereumDao), repo.NewNormalAccountRepo(ethereumDao))
		accountMngWriter := ethm.NewRetryProcess("account", accountMng, errorRepo)
		ethWriterCtl.RegisterTxWriter(accountMngWriter)
		svr.Add(accountMngWriter)
	}

	if conf.SaveAccount {
		ethereumDao.InitContractAccountTb()
		ethereumDao.InitEthAccountTb()
		contractMng := ethm.NewContractManager(ethm.NewEthRpcExecutor(conf.RpcUrl), repo.NewContractRepo(ethereumDao))
		contractMngWriter := ethm.NewRetryProcess("contract", contractMng, errorRepo)
		ethWriterCtl.RegisterTxWriter(contractMngWriter)
		svr.Add(contractMngWriter)
	}

	if conf.SaveTransaction {
		var (
			contractTxDb      = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthContractTransactionDbConfig()), "contract_transaction")
			transactionDb     = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthTransactionDbConfig()), "transaction")
			normalTxDbCache   = dao.NewDbCache("err_data/normal_traction.sql", conf.WriteToDbInterval, transactionDb)
			contractTxDbCache = dao.NewDbCache("err_data/contract_traction.sql", conf.WriteToDbInterval, contractTxDb)
			normalTranDao     = dao.NewEthereumTransactionDao(transactionDb, contractTxDb, normalTxDbCache, contractTxDbCache)
		)

		var (
			transactionRepoFactory = repo.NewTransactionRepositoryFactory(conf.SaveType, conf.WriteToDbInterval, "err_data/transaction_data.sql", normalTranDao)
			transactionRepo        = transactionRepoFactory.NewTransactionRepository()
			transactionWriter      = ethm.NewTransactionWriter(ethm.NewEthRpcExecutor(conf.RpcUrl), transactionRepo)
			transactionReWriter    = ethm.NewRetryProcess("transaction", transactionWriter, errorRepo)
		)
		transactionWriter.SetWithBalance(conf.WithTxBalance)
		ethWriterCtl.RegisterTxWriter(transactionReWriter)
		svr.Add(transactionReWriter)
		switch conf.SaveType {
		case repo.SaveTypeInHive, repo.SaveTypeInSqlFile:
			r, ok := transactionRepo.(runner.Runner)
			if ok {
				svr.Add(r)
			}
		default:
			svr.Add(normalTxDbCache)
			svr.Add(contractTxDbCache)
		}
	}

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
