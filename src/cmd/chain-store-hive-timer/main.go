package main

//
//import (
//	"context"
//	"flag"
//	"fmt"
//	"github.com/ville-vv/eth-chain-store/src/common/conf"
//	"github.com/ville-vv/eth-chain-store/src/common/go_exec"
//	"github.com/ville-vv/eth-chain-store/src/common/hive"
//	"github.com/ville-vv/eth-chain-store/src/common/log"
//	"github.com/ville-vv/eth-chain-store/src/domain/async"
//	"github.com/ville-vv/eth-chain-store/src/domain/repo"
//	"github.com/ville-vv/eth-chain-store/src/infra/dao"
//	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
//	"github.com/ville-vv/eth-chain-store/src/infra/monitor"
//	"github.com/ville-vv/eth-chain-store/src/server"
//	"github.com/ville-vv/vilgo/vlog"
//	"github.com/ville-vv/vilgo/vstore"
//	"os"
//	"runtime"
//	"time"
//)
//
//var (
//	//syncInterval      string
//	//fastSyncInterval  string
//	rpcEndpoint       string
//	dbUser            string
//	dbPassword        string
//	dbHost            string
//	dbPort            string
//	logFile           string
//	debug             bool
//	maxPullNum        int
//	maxWriteNum       int
//	isMaxProcs        bool
//	isHelp            bool
//	writeToDbInterval int
//)
//
//func cmdFlagParse() {
//	//flag.StringVar(&syncInterval, "si", "15", "the interval to sync latest block number")
//	//flag.StringVar(&fastSyncInterval, "fsi", "1000", "the interval fast to sync the block number  before  the latest ms")
//	flag.StringVar(&rpcEndpoint, "rpc_url", "", "eth rpc endpoint")
//	flag.StringVar(&dbUser, "db_user", "", "the database user")
//	flag.StringVar(&dbPassword, "db_passwd", "", "the database password")
//	flag.StringVar(&dbHost, "db_host", "", "the database host")
//	flag.StringVar(&dbPort, "db_port", "", "the database port")
//	flag.StringVar(&logFile, "logFile", "", "the log file path and file name")
//	//flag.IntVar(&maxPullNum, "max_pull_num", 1, "the max thread number for sync block information from chain")
//	//flag.IntVar(&maxWriteNum, "max_write_num", 5, "the max thread number for write block information to db")
//	flag.IntVar(&writeToDbInterval, "wi", 2, "the interval time that write to mysql from memory n/s")
//	flag.BoolVar(&debug, "debug", false, "open debug logs")
//	flag.BoolVar(&isHelp, "help", false, "help")
//	flag.BoolVar(&isMaxProcs, "max_procs", false, "the max process core the value is true or false")
//	flag.Parse()
//	fmt.Println(rpcEndpoint, logFile)
//	if isHelp {
//		flag.PrintDefaults()
//		os.Exit(1)
//	}
//	if rpcEndpoint == "" {
//		fmt.Println("rpc_endpoint is empty")
//		flag.PrintDefaults()
//		os.Exit(-1)
//	}
//}
//
//func buildService() go_exec.Runner {
//	//dbMysql := vstore.MakeDb(conf.GetEthereumHiveMapDbConfig())
//	//hiveCli, err := hive.New(conf.GetHiveEthereumDb())
//	//if err != nil {
//	//	panic(err)
//	//}
//	//
//	//var (
//	//	ethereumCacheDb = dao.NewDbCache("err_data/ethereum_map_hive_01.sql", writeToDbInterval, dbMysql)
//	//	ethereumDao     = dao.NewEthereumDao(dbMysql, ethereumCacheDb)
//	//
//	//	cursorRepo           = dao.NewEthereumMapHive("err_data/ethereum_map_hive_02.sql", dbMysql, hiveCli, writeToDbInterval)
//	//	errorDao             = dao.NewSyncErrorDao(dbMysql)
//	//	errorRepo            = repo.NewSyncErrorRepository(errorDao)
//	//	latestBlockNumGetter = async.NewLatestBlockNumberCache(ethrpc.NewClient(rpcEndpoint))
//	//)
//	//
//	//var (
//	//	contractDataCursor = async.NewDataCursorAggregate(async.CursorTypeContractTx, cursorRepo)
//	//	contractAccount    = async.NewContractAccountService(ethrpc.NewClient(rpcEndpoint), repo.NewContractAccountRepo(ethereumDao))
//	//	contractAccountPcr = async.NewDataProcessorCtl(contractAccount, contractDataCursor, errorRepo, latestBlockNumGetter)
//	//)
//	//
//	//var (
//	//	ethDataCursor = async.NewDataCursorAggregate(async.CursorTypeEthereumTx, cursorRepo)
//	//	ethAccount    = async.NewEthAccountService(ethrpc.NewClient(rpcEndpoint), repo.NewNormalAccountRepo(ethereumDao))
//	//	ethAccountPcr = async.NewDataProcessorCtl(ethAccount, ethDataCursor, errorRepo, latestBlockNumGetter)
//	//)
//	//contractAccountPcr.SetName("合约账户扫描")
//	//ethAccountPcr.SetName("以太坊账户扫描")
//	//
//	//timerSvr := server.NewTimerServer()
//	//timerSvr.Add(ethereumCacheDb, latestBlockNumGetter, cursorRepo, contractAccountPcr, ethAccountPcr)
//	//
//	//return timerSvr
//}
//
//func main() {
//	cmdFlagParse()
//	if isMaxProcs {
//		runtime.GOMAXPROCS(runtime.NumCPU())
//	}
//	log.Init()
//	go_exec.Go(monitor.StartMonitor)
//	ctx, cancel := context.WithCancel(context.Background())
//	go_exec.Go(func() {
//		select {
//		case <-conf.GlobalProgramFinishSigmal:
//			cancel()
//			return
//		}
//	})
//
//	go_exec.Run(ctx, buildService())
//	close(conf.GlobalExitSignal)
//	time.Sleep(time.Second * 5)
//	vlog.INFO("定时扫描工具退出")
//}

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

func main() {
	var echoTimes int
	var sbc string
	var cmdPrint = &cobra.Command{
		Use:   "print [string to print]",
		Short: "Print anything to the screen",
		Long: `print is for printing anything back to the screen.
For many years people have printed back to the screen.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Print: " + strings.Join(args, " "))
			fmt.Println(sbc)
		},
	}

	cmdPrint.Flags().StringVar(&sbc, "sbc", "sbc", "xxdg")

	var cmdEcho = &cobra.Command{
		Use:   "echo [string to echo]",
		Short: "Echo anything to the screen",
		Long: `echo is for echoing anything back.
Echo works a lot like print, except it has a child command.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Echo: " + strings.Join(args, " "))
		},
	}

	var cmdTimes = &cobra.Command{
		Use:   "times [string to echo]",
		Short: "Echo anything to the screen more times",
		Long: `echo things multiple times back to the user by providing
a count and a string.`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for i := 0; i < echoTimes; i++ {
				fmt.Println("Echo: " + strings.Join(args, " "))
			}

			cmd.Flag("xxx")
		},
	}

	cmdTimes.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdPrint, cmdEcho)
	cmdEcho.AddCommand(cmdTimes)
	rootCmd.Execute()
}
