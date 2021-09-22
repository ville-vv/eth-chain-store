package async

import (
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/domain/ethm"
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/cache"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
)

// ContractAccountService
type ContractAccountService struct {
	rpcCli      ethrpc.EthRPC
	accountRepo repo.ContractAccountRepository
	accountMng  *ethm.ContractAccountManager
}

func NewContractAccountService(rpcCli ethrpc.EthRPC, accountRepo repo.ContractAccountRepository) *ContractAccountService {
	return &ContractAccountService{
		rpcCli:      rpcCli,
		accountRepo: accountRepo,
		accountMng: ethm.NewContractAccountManager(
			rpcCli, accountRepo, ethm.NewRingStrListV2(), cache.NewRingCache(),
		)}
}

func (sel *ContractAccountService) Process(latestNumber int64, record *model.TransactionRecord) error {
	return sel.accountMng.UpdateAccount(&model.TransactionData{
		LatestNumber:    fmt.Sprintf("%d", latestNumber),
		ContractAddress: record.ContractAddress,
		TimeStamp:       record.TxTime,
		BlockHash:       record.BlockHash,
		BlockNumber:     record.BlockNumber,
		From:            record.FromAddr,
		Hash:            record.TxHash,
		To:              record.ToAddr,
	})
}
