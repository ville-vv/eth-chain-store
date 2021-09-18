package async

import (
	"github.com/ville-vv/eth-chain-store/src/domain/entity"
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
)

type TxDataGetter interface {
	GetTxData() ([]*model.TransactionRecord, error)
	Finish() error
}

type ContractExistChecker interface {
	Exist(key string) bool
}

type ContractService struct {
	rpcCli       ethrpc.EthRPC
	dataCursor   TxDataGetter
	contractRepo repo.ContractRepository
	errorRepo    repo.SyncErrorRepository
}

func (sel *ContractService) Process(data *model.TransactionRecord) error {
	return sel.process(data.ContractAddress, data.TxTime)
}

func (sel *ContractService) process(address string, timestamp string) error {
	// 判断是数据否存在
	if sel.contractRepo.IsContractExist(address) {
		return nil
	}
	contract := entity.NewContract(sel.rpcCli, sel.contractRepo)
	contract.SetAddress(address)
	contract.SetPublishTime(timestamp)
	if err := contract.SetErc20ContentFromRpc(); err != nil {
		return err
	}
	// 数据不存在就创建
	return contract.CreateRecord()
}
