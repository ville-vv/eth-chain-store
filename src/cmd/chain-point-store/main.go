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
	"os"
	"runtime"
	"time"
)

var (
	syncInterval      string
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

	startBlockNumber int // 开始区块
	endBlockNumber   int // 结束区块
)

var (
	errDataFile    *os.File
	syncConfigFile *os.File
)

func init() {
	var err error
	errDataFile, err = os.OpenFile("err_data/err_data_log", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	syncConfigFile, err = os.OpenFile("sync_data.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
}

func cmdFlagParse() {
	flag.StringVar(&rpcEndpoint, "rpc_endpoint", "http://localhost:8545", "eth rpc endpoint")
	flag.StringVar(&dbUser, "db_user", "", "the database user")
	flag.StringVar(&dbPassword, "db_passwd", "", "the database password")
	flag.StringVar(&dbHost, "db_host", "", "the database host")
	flag.StringVar(&dbPort, "db_port", "", "the database port")
	flag.StringVar(&logFile, "logFile", "", "the log file path and file name")
	flag.IntVar(&maxPullNum, "max_sync_thread", 1, "the max thread number for sync block information from chain")
	flag.IntVar(&maxWriteNum, "max_write_thread", 5, "the max thread number for write block information to db")
	flag.IntVar(&writeToDbInterval, "wi", 2, "the interval that write to file ")

	flag.IntVar(&startBlockNumber, "start_number", 0, "the start block number need to sync")
	flag.IntVar(&endBlockNumber, "end_number", 0, "the end block number need to sync ")

	flag.BoolVar(&debug, "debug", false, "open debug logs")
	flag.BoolVar(&isHelp, "help", false, "help")
	flag.BoolVar(&isMaxProcs, "max_cpu", false, "use the max cpu process numbers")
	flag.Parse()
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

	var (
		txDataRepo = dao.NewHiveDataFile("hive_data/transaction_record.txt", "hive_data/contract_transaction_record.txt")

		transactionWriter   = ethm.NewTransactionWriter(ethm.NewEthRpcExecutor(rpcEndpoint), txDataRepo)
		transactionReWriter = ethm.NewRetryProcess("transaction", maxWriteNum, transactionWriter, repo.NewSyncErrorRepositoryV2(errDataFile))

		mqPublish       = mqp.NewMDP(maxWriteNum*2, vlog.INFO)
		txWriterPublish = ethm.NewEthereumPublisher(mqPublish)

		ethMng = ethm.NewEthereumManager(ethm.NewEthRpcExecutor(rpcEndpoint), txWriterPublish)

		syncControl = ethm.NewSyncBlockControlWithOpt(
			&ethm.OptionConfig{
				StartBlockNumber: 0,
				EndBlockNumber:   0,
				MaxSyncThreads:   maxPullNum,
				EthRpcCli:        ethm.NewEthRpcExecutor(rpcEndpoint),
				BknRepo:          repo.NewBlockNumberRepoV2(syncConfigFile),
			},
		)

		serviceRun = server.NewSyncBlockChainServiceV2(ethMng, syncControl)
	)

	mqPublish.SubScribe(transactionReWriter)
	transactionReWriter.SetMonitor(&monitor.TxWriteProcessNum)

	svr := &server.Server{}
	svr.Add(serviceRun, mqPublish)
	svr.Add(transactionReWriter, txDataRepo)

	return svr
}

func main() {
	cmdFlagParse()
	if isMaxProcs {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	log.Init()
	runner.Go(monitor.StartMonitor)
	runner.Run(buildService())
	close(conf.GlobalExitSignal)
	time.Sleep(time.Second * 5)
}
