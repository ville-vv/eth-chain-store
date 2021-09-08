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
	"github.com/ville-vv/vilgo/vfile"
	"os"
	"path"
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

	startBlockNumber int64 // 开始区块
	endBlockNumber   int64 // 结束区块
)

var (
	errDataFile    *os.File
	syncConfigFile *os.File
)

func openFile(fileName string, flag int) (*os.File, error) {
	dirPath := path.Dir(fileName)
	if !vfile.PathExists(dirPath) {
		err := os.Mkdir(dirPath, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	f, err := os.OpenFile(fileName, flag, 0644)
	if err != nil {
		return nil, err
	}
	return f, err
}

func Init() {
	log.Init()
	var err error
	errDataFile, err = openFile("err_data/err_data_log", os.O_WRONLY|os.O_CREATE)
	if err != nil {
		panic(err)
	}

	syncConfigFile, err = openFile("sync_data.json", os.O_RDWR|os.O_CREATE)
	if err != nil {
		panic(err)
	}
}

func exit() {
	_ = errDataFile.Close()
	_ = syncConfigFile.Close()
}

func cmdFlagParse() {
	flag.StringVar(&rpcEndpoint, "rpc_url", "http://localhost:8545", "eth rpc url")
	flag.StringVar(&logFile, "logFile", "", "the log file path and file name")
	flag.IntVar(&maxPullNum, "max_sync_thread", 1, "the max thread number for sync block information from chain")
	flag.IntVar(&maxWriteNum, "max_write_thread", 5, "the max thread number for write block information to db")

	flag.Int64Var(&startBlockNumber, "start_number", 0, "the start block number need to sync")
	flag.Int64Var(&endBlockNumber, "end_number", 0, "the end block number need to sync ")

	flag.BoolVar(&debug, "debug", false, "open debug logs")
	flag.BoolVar(&isHelp, "help", false, "help")
	flag.BoolVar(&isMaxProcs, "max_cpu", false, "use the max cpu process numbers")
	flag.IntVar(&writeToDbInterval, "wi", 2, "the max thread number for write block information to db")
	flag.Parse()
	if isHelp {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if rpcEndpoint == "" {
		fmt.Println("rpc_url is empty")
		flag.PrintDefaults()
		os.Exit(-1)
	}
}

func buildService() go_exec.Runner {
	//debug = true
	//if debug {
	//	rpcEndpoint = "http://172.16.16.115:8545"
	//	startBlockNumber = 5000000
	//	endBlockNumber = 5001000
	//	maxPullNum = 10
	//	maxWriteNum = 100
	//}

	var (
		txDataRepo = dao.NewHiveDataFile(
			"hive_data/transaction_record.txt",
			"hive_data/contract_transaction_record.txt",
			int64(writeToDbInterval))

		transactionWriter   = ethm.NewTransactionWriter(ethm.NewEthRpcExecutor(rpcEndpoint), txDataRepo)
		transactionReWriter = ethm.NewRetryProcess("transaction", transactionWriter, repo.NewSyncErrorRepositoryV2(errDataFile))

		ethDataWriter = ethm.NewEthWriterControl(maxPullNum, transactionReWriter)
		ethMng        = ethm.NewEthereumManager(ethm.NewEthRpcExecutor(rpcEndpoint), ethDataWriter)

		syncControl = ethm.NewSyncBlockControlWithOpt(
			&ethm.OptionConfig{
				StartBlockNumber: startBlockNumber,
				EndBlockNumber:   endBlockNumber,
				MaxSyncThreads:   maxPullNum,
				EthRpcCli:        ethm.NewEthRpcExecutor(rpcEndpoint),
				BknRepo:          repo.NewBlockNumberRepoV2(syncConfigFile),
			},
		)
		serviceRun = server.NewSyncBlockChainServiceV2(ethMng, syncControl)
	)

	svr := &server.Server{}
	svr.Add(serviceRun, ethDataWriter)
	svr.Add(transactionReWriter, txDataRepo)

	return svr
}

func main() {
	cmdFlagParse()
	if isMaxProcs {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	Init()
	runner.Go(monitor.StartMonitor)

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
	exit()
}
