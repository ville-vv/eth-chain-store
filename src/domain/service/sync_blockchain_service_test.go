package service

import (
	"context"
	"github.com/ville-vv/eth-chain-store/src/common/conf"
	"github.com/ville-vv/eth-chain-store/src/common/log"
	"github.com/ville-vv/eth-chain-store/src/domain/ethm"
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/dao"
	"github.com/ville-vv/eth-chain-store/src/infra/mqp"
	"github.com/ville-vv/vilgo/runner"
	"github.com/ville-vv/vilgo/vlog"
	"github.com/ville-vv/vilgo/vstore"
	"testing"
)

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

func TestNewSyncBlockChainService(t *testing.T) {
	log.Init()
	var rpcEndpoint = "http://172.16.16.115:8545"
	var maxWriteNum = 2
	var maxPullNum = 1

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
		contractMng         = ethm.NewContractManager(ethm.NewEthRpcExecutor(rpcEndpoint, ""), repo.NewContractRepo(ethereumDao))
		accountMng          = ethm.NewAccountManager(ethm.NewEthRpcExecutor(rpcEndpoint, ""), repo.NewContractAccountRepo(ethereumDao), repo.NewNormalAccountRepo(ethereumDao))
		transactionWriter   = ethm.NewTransactionWriter(ethm.NewEthRpcExecutor(rpcEndpoint, ""), repo.NewTransactionRepo(normalTranDao))
		accountMngWriter    = ethm.NewEthRetryWriter("account", maxWriteNum, accountMng, repo.NewSyncErrorRepository(errorDao))
		contractMngWriter   = ethm.NewEthRetryWriter("contract", maxWriteNum, contractMng, repo.NewSyncErrorRepository(errorDao))
		transactionReWriter = ethm.NewEthRetryWriter("transaction", maxWriteNum, transactionWriter, repo.NewSyncErrorRepository(errorDao))

		mqPublish       = mqp.NewMDP(vlog.INFO)
		txWriterPublish = ethm.NewEthereumPublisher(mqPublish)

		ethMng       = ethm.NewEthereumManager(ethm.NewEthRpcExecutor(rpcEndpoint, ""), txWriterPublish)
		bkNumCounter = ethm.NewSyncBlockNumberCounterV2(maxPullNum, ethm.NewEthRpcExecutor(rpcEndpoint, ""), repo.NewBlockNumberRepo(ethBlockNumberDao))
		serviceRun   = NewSyncBlockChainServiceV2(ethMng, bkNumCounter)
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
	runner.Run(svr)

}
