package async

import (
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/domain/ethm"
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
)

// EthAccountService
type EthAccountService struct {
	rpcCli      ethrpc.EthRPC
	accountRepo repo.EthAccountRepository
	accountMng  *ethm.NormalAccountManager
}

func NewEthAccountService(rpcCli ethrpc.EthRPC, accountRepo repo.EthAccountRepository) *EthAccountService {
	return &EthAccountService{
		rpcCli:      rpcCli,
		accountRepo: accountRepo,
		accountMng:  ethm.NewNormalAccountManager(rpcCli, accountRepo, ethm.NewRingStrListV2())}
}

func (sel *EthAccountService) Process(latestNumber int64, record *model.TransactionRecord) error {
	return sel.accountMng.UpdateAccount(&model.TransactionData{
		LatestNumber: fmt.Sprintf("%d", latestNumber),
		TimeStamp:    record.TxTime,
		BlockHash:    record.BlockHash,
		BlockNumber:  record.BlockNumber,
		From:         record.FromAddr,
		Hash:         record.TxHash,
		To:           record.ToAddr,
	})
}
