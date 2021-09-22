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
}

func (sel *ContractService) Process(data *model.TransactionRecord) error {
	contract := entity.NewContract(sel.rpcCli, sel.contractRepo)
	contract.SetAddress(data.ContractAddress)
	contract.SetPublishTime(data.TxTime)
	// 创建合约信息
	if err := contract.CreateRecord(); err != nil {
		return err
	}
	return nil
}
