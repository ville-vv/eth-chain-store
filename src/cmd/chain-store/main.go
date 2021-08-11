package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/common/conf"
	"github.com/ville-vv/eth-chain-store/src/common/log"
	"github.com/ville-vv/eth-chain-store/src/domain/ethm"
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/domain/service"
	"github.com/ville-vv/eth-chain-store/src/infra/dao"
	"github.com/ville-vv/eth-chain-store/src/infra/mqp"
	"github.com/ville-vv/vilgo/runner"
	"github.com/ville-vv/vilgo/vlog"
	"github.com/ville-vv/vilgo/vstore"
	"os"
	"runtime"
	"time"
)

var (
	syncInterval     string
	fastSyncInterval string
	rpcEndpoint      string
	dbUser           string
	dbPassword       string
	dbHost           string
	dbPort           string
	logFile          string
	debug            bool
	maxPullNum       int
	maxWriteNum      int
	isMaxProcs       bool
	isHelp           bool
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
	var (
		//businessDb = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthBusinessDbConfig()), "business")
		ethereumDb = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthereumDbConfig()), "ethereum")
		//contractDb    = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthContractDbConfig()), "contract")
		contractTxDb  = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthContractTransactionDbConfig()), "contract_transaction")
		transactionDb = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthTransactionDbConfig()), "transaction")
	)

	var (
		normalTranDao = dao.NewEthereumTransactionDao(transactionDb, contractTxDb)
		//erc20TokenCfgDao  = dao.NewErc20TokenConfigDao(businessDb)
		ethereumDao       = dao.NewEthereumDao(ethereumDb)
		ethBlockNumberDao = dao.NewEthereumBlockNumberDao(ethereumDb)
		errorDao          = dao.NewSyncErrorDao(ethereumDb)
	)
	var (
		// https://mainnet.infura.io/v3/ecc309a045134205b5c2b58481d7923d
		// https://mainnet.infura.io/v3/21628f8f9b9b423a9ea05a708016b119
		contractMng         = ethm.NewContractManager(ethm.NewEthRpcExecutor(rpcEndpoint, ""), repo.NewContractRepo(ethereumDao))
		accountMng          = ethm.NewAccountManager(ethm.NewEthRpcExecutor(rpcEndpoint, ""), repo.NewContractAccountRepo(ethereumDao), repo.NewNormalAccountRepo(ethereumDao))
		transactionWriter   = ethm.NewTransactionWriter(ethm.NewEthRpcExecutor(rpcEndpoint, ""), repo.NewTransactionRepo(normalTranDao))
		accountMngWriter    = ethm.NewEthRetryWriter("account", maxWriteNum, accountMng, repo.NewSyncErrorRepository(errorDao))
		contractMngWriter   = ethm.NewEthRetryWriter("contract", maxWriteNum, contractMng, repo.NewSyncErrorRepository(errorDao))
		transactionReWriter = ethm.NewEthRetryWriter("transaction", maxWriteNum, transactionWriter, repo.NewSyncErrorRepository(errorDao))
		//txWriter            = ethm.NewEthereumWriter(filter, accountMngWriter, contractMngWriter, transactionReWriter)

		mqPublish       = mqp.NewMDP(vlog.ERROR)
		txWriterPublish = ethm.NewEthereumPublisher(mqPublish)
		//serviceRun      = service.NewSyncBlockChainService(maxPullNum, ethm.NewEthRpcExecutor(rpcEndpoint, ""), txWriterPublish, repo.NewBlockNumberRepo(ethBlockNumberDao))

		ethMng       = ethm.NewEthereumManager(ethm.NewEthRpcExecutor(rpcEndpoint, ""), txWriterPublish)
		bkNumCounter = ethm.NewSyncBlockNumberCounterV2(maxPullNum, ethm.NewEthRpcExecutor(rpcEndpoint, ""), repo.NewBlockNumberRepo(ethBlockNumberDao))
		serviceRun   = service.NewSyncBlockChainServiceV2(ethMng, bkNumCounter)
	)
	mqPublish.SubScribe(accountMngWriter)
	mqPublish.SubScribe(contractMngWriter)
	mqPublish.SubScribe(transactionReWriter)

	svr := &Server{}
	svr.Add(serviceRun)
	svr.Add(mqPublish)
	svr.Add(accountMngWriter)
	svr.Add(contractMngWriter)
	svr.Add(transactionReWriter)

	return svr
}

type Server struct {
	runner []runner.Runner
}

func (s *Server) Scheme() string {
	return "Server"
}

func (s *Server) Init() error {
	for _, r := range s.runner {
		return r.Init()
	}
	return nil
}

func (s *Server) Start() error {
	for _, r := range s.runner {
		runner.Go(func() {
			_ = r.Start()
		})
	}
	return nil
}

func (s *Server) Exit(ctx context.Context) error {
	for _, r := range s.runner {
		_ = r.Exit(ctx)
	}
	return nil
}

func (s *Server) Add(r runner.Runner) {
	s.runner = append(s.runner, r)
}

func main() {
	cmdFlagParse()
	if isMaxProcs {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	log.Init() //12900839
	runner.Run(buildService())
	time.Sleep(time.Second * 5)
}
