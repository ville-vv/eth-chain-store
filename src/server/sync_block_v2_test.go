package server

import (
	"context"
	"github.com/ville-vv/eth-chain-store/src/common/conf"
	"github.com/ville-vv/eth-chain-store/src/common/go_exec"
	"github.com/ville-vv/eth-chain-store/src/common/log"
	"github.com/ville-vv/eth-chain-store/src/domain/ethm"
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/dao"
	"github.com/ville-vv/vilgo/vstore"
	"testing"
)

func TestNewSyncBlockChainService(t *testing.T) {
	log.Init()
	var rpcEndpoint = "http://172.16.16.115:8545"
	var maxWriteNum = 10
	var maxPullNum = 1

	rpcEndpoint = "https://mainnet.infura.io/v3/21628f8f9b9b423a9ea05a708016b119"

	var (
		//businessDb = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthBusinessDbConfig()), "business")
		ethereumDb = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthereumDbConfig()), "ethereum")
		//contractDb    = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthContractDbConfig()), "contract")
		contractTxDb  = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthContractTransactionDbConfig()), "contract_transaction")
		transactionDb = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthTransactionDbConfig()), "transaction")
	)

	var (
		normalTxDbCache   = dao.NewDbCache("a", 1, transactionDb)
		contractTxDbCache = dao.NewDbCache("b", 1, contractTxDb)
		normalTranDao     = dao.NewEthereumTransactionDao(transactionDb, contractTxDb, normalTxDbCache, contractTxDbCache)
		//erc20TokenCfgDao  = dao.NewErc20TokenConfigDao(businessDb)
		ethereumCacheDb   = dao.NewDbCache("c", 1, ethereumDb)
		ethereumDao       = dao.NewEthereumDao(ethereumDb, ethereumCacheDb)
		ethBlockNumberDao = dao.NewEthereumBlockNumberDao(ethereumDb)
		errorDao          = dao.NewSyncErrorDao(ethereumDb)
	)
	var (
		contractMng         = ethm.NewContractManager(ethm.NewEthRpcExecutor(rpcEndpoint), repo.NewContractRepo(ethereumDao))
		accountMng          = ethm.NewAccountManager(ethm.NewEthRpcExecutor(rpcEndpoint), repo.NewContractAccountRepo(ethereumDao), repo.NewNormalAccountRepo(ethereumDao))
		transactionWriter   = ethm.NewTransactionWriter(ethm.NewEthRpcExecutor(rpcEndpoint), repo.NewTransactionRepo(normalTranDao))
		accountMngWriter    = ethm.NewRetryProcess("account", accountMng, repo.NewSyncErrorRepository(errorDao))
		contractMngWriter   = ethm.NewRetryProcess("contract", contractMng, repo.NewSyncErrorRepository(errorDao))
		transactionReWriter = ethm.NewRetryProcess("transaction", transactionWriter, repo.NewSyncErrorRepository(errorDao))

		txWriterPublish = ethm.NewEthWriterControl(maxWriteNum, accountMngWriter, contractMngWriter, transactionReWriter)

		ethMng       = ethm.NewEthereumManager(ethm.NewEthRpcExecutor(rpcEndpoint), txWriterPublish)
		bkNumCounter = ethm.NewSyncBlockControl(maxPullNum, ethm.NewEthRpcExecutor(rpcEndpoint), repo.NewBlockNumberRepo(ethBlockNumberDao))
		serviceRun   = NewSyncBlockChainServiceV2(ethMng, bkNumCounter)
	)

	svr := &Server{}
	svr.Add(serviceRun)

	svr.Add(accountMngWriter)
	svr.Add(contractMngWriter)
	svr.Add(transactionReWriter)
	svr.Add(ethereumCacheDb)
	go_exec.Run(context.Background(), svr)

}
