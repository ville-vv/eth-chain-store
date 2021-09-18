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
	"github.com/ville-vv/vilgo/vstore"
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
	var (
		ethereumHiveMapDb = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthereumHiveMapDbConfig()), "ethereum_hive_map")
		ethBlockNumberDao = dao.NewEthereumBlockNumberDao(ethereumHiveMapDb)

		txDataRepo = dao.NewTransactionHiveDao("err_data/transaction_hive.sql", conf.GetHiveEthereumDb())
		errorRepo  = repo.NewSyncErrorRepository(dao.NewSyncErrorDao(ethereumHiveMapDb))

		transactionWriter   = ethm.NewTransactionWriter(ethm.NewEthRpcExecutor(rpcEndpoint), txDataRepo)
		transactionReWriter = ethm.NewRetryProcess("transaction", transactionWriter, errorRepo)

		ethDataWriter = ethm.NewEthWriterControl(maxWriteNum, transactionReWriter)
		ethMng        = ethm.NewEthereumManager(ethm.NewEthRpcExecutor(rpcEndpoint), ethDataWriter)
		syncControl   = ethm.NewSyncBlockControlWithOpt(
			&ethm.OptionConfig{
				StartBlockNumber: startBlockNumber,
				EndBlockNumber:   endBlockNumber,
				MaxSyncThreads:   maxPullNum,
				EthRpcCli:        ethm.NewEthRpcExecutor(rpcEndpoint),
				BknRepo:          repo.NewBlockNumberRepo(ethBlockNumberDao),
			},
		)
	)
	syncControl.SetPuller(ethMng)
	svr := &server.Server{}
	svr.Add(syncControl, ethDataWriter)
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
}
