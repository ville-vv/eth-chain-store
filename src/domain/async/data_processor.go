package async

import (
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
)

type DataProcessor interface {
	Process(latestNumber int64, record *model.TransactionRecord) error
}

// 数据处理控制器
type DataProcessorCtl struct {
	rpcCli     ethrpc.EthRPC
	processor  DataProcessor
	dataCursor TxDataGetter
	errorRepo  repo.SyncErrorRepository
}

func NewDataProcessorCtl(processor DataProcessor, dataCursor TxDataGetter, errorRepo repo.SyncErrorRepository) *DataProcessorCtl {
	return &DataProcessorCtl{processor: processor, dataCursor: dataCursor, errorRepo: errorRepo}
}

func (sel *DataProcessorCtl) Process() error {
	// 获取最新区块
	latestNumber, err := sel.rpcCli.GetBlockNumber()
	if err != nil {

		return err
	}
	// 1.获取交易数据
	dataList, err := sel.dataCursor.GetTxData()
	if err != nil {
		return err
	}

	for _, data := range dataList {
		if err := sel.processor.Process(int64(latestNumber), data); err != nil {
			// 记录错误的数据
			_ = sel.errorRepo.WriterErrorRecord("", data.TxHash, data.BlockNumber, err)
		}
	}
	//
	_ = sel.dataCursor.Finish()
	return nil
}
