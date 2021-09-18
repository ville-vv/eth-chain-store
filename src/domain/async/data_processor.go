package async

import (
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
)

type DataProcessor interface {
	Process(record *model.TransactionRecord) error
}

// 数据处理控制器
type DataProcessorCtl struct {
	processor  DataProcessor
	dataCursor TxDataGetter
	errorRepo  repo.SyncErrorRepository
}

func (sel *DataProcessorCtl) Process() error {
	// 1.获取交易数据
	dataList, err := sel.dataCursor.GetTxData()
	if err != nil {
		return err
	}
	// 数据去重
	for _, data := range dataList {
		if err := sel.processor.Process(data); err != nil {
			// 记录错误的数据
			_ = sel.errorRepo.WriterErrorRecord("", data.TxHash, data.BlockNumber, err)
		}
	}
	//
	_ = sel.dataCursor.Finish()
	return nil
}
