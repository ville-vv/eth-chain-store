package service

import (
	"github.com/ville-vv/eth-chain-store/src/common/conf"
	"github.com/ville-vv/eth-chain-store/src/common/log"
	"github.com/ville-vv/eth-chain-store/src/domain/ethm"
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/dao"
	"github.com/ville-vv/vilgo/vstore"
	"testing"
)

func TestNewSyncBlockChainService(t *testing.T) {
	log.Init()
	var (
		businessDb = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthBusinessDbConfig()), "business")
		ethereumDb = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthereumDbConfig()), "ethereum")
		//contractDb    = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthContractDbConfig()), "contract")
		contractTxDb  = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthContractTransactionDbConfig()), "contract_transaction")
		transactionDb = dao.NewMysqlDB(vstore.MakeDb(conf.GetEthTransactionDbConfig()), "transaction")
	)

	var (
		normalTranDao     = dao.NewEthereumTransactionDao(transactionDb, contractTxDb)
		erc20TokenCfgDao  = dao.NewErc20TokenConfigDao(businessDb)
		ethereumDao       = dao.NewEthereumDao(ethereumDb)
		ethBlockNumberDao = dao.NewEthereumBlockNumberDao(ethereumDb)
	)
	var (
		// https://mainnet.infura.io/v3/ecc309a045134205b5c2b58481d7923d
		// https://mainnet.infura.io/v3/21628f8f9b9b423a9ea05a708016b119
		ehtrpcCli         = ethm.NewEthRpcExecutor("https://mainnet.infura.io/v3/ecc309a045134205b5c2b58481d7923d", "")
		filter            = ethm.NewEthereumWriteFilter(erc20TokenCfgDao)
		contractMng       = ethm.NewContractManager(ehtrpcCli, repo.NewContractRepo(ethereumDao))
		accountMng        = ethm.NewAccountManager(ehtrpcCli, repo.NewContractAccountRepo(ethereumDao), repo.NewNormalAccountRepo(ethereumDao))
		transactionWriter = ethm.NewTransactionWriter(ehtrpcCli, repo.NewTransactionRepo(normalTranDao))
		txWriter          = ethm.NewEthereumWriter(filter, accountMng, contractMng, transactionWriter)
	)
	syncSvc := NewSyncBlockChainService(1, ehtrpcCli, txWriter, repo.NewBlockNumberRepo(ethBlockNumberDao))
	if err := syncSvc.Start(); err != nil {
		t.Error(err)
		return
	}
}
