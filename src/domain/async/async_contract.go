package async

import (
	"github.com/ville-vv/eth-chain-store/src/domain/entity"
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
)

type ContractExistChecker interface {
	Exist(key string) bool
}

type ContractService struct {
	rpcCli       ethrpc.EthRPC
	contractRepo repo.ContractRepository
}

func NewContractService(rpcCli ethrpc.EthRPC, contractRepo repo.ContractRepository) *ContractService {
	return &ContractService{rpcCli: rpcCli, contractRepo: contractRepo}
}

func (sel *ContractService) Process(latestNum int64, data *model.TransactionRecord) error {
	contract := entity.NewContract(sel.rpcCli, sel.contractRepo)
	contract.SetAddress(data.ContractAddress)
	contract.SetPublishTime(data.TxTime)
	// 创建合约信息
	if err := contract.CreateRecord(); err != nil {
		return err
	}
	return nil
}
