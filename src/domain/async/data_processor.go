package async

import (
	"context"
	"github.com/ville-vv/eth-chain-store/src/common/go_exec"
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/vilgo/vlog"
	"time"
)

type TxDataGetter interface {
	GetTxData() ([]*model.TransactionRecord, error)
	Finish() error
}

type DataProcessor interface {
	Process(latestNumber int64, record *model.TransactionRecord) error
}

type BlockNumberGetter interface {
	GetBlockNumber() (uint64, error)
}

// 数据处理控制器
type DataProcessorCtl struct {
	processorList []DataProcessor
	dataCursor    TxDataGetter
	blockNum      BlockNumberGetter
	errorRepo     repo.SyncErrorRepository
	isStop        bool
	waitFinish    chan int
	name          string
}

func NewDataProcessorCtl(dataCursor TxDataGetter, errorRepo repo.SyncErrorRepository, blockNum BlockNumberGetter) *DataProcessorCtl {
	return &DataProcessorCtl{dataCursor: dataCursor, errorRepo: errorRepo, blockNum: blockNum, waitFinish: make(chan int)}
}

func (sel *DataProcessorCtl) AddProcess(prcList ...DataProcessor) {
	sel.processorList = append(sel.processorList, prcList...)
}

func (sel *DataProcessorCtl) SetName(name string) {
	sel.name = name
}

func (sel *DataProcessorCtl) Start() error {
	go_exec.Go(func() {
		for {
			time.Sleep(time.Second)
			if sel.isStop {
				<-sel.waitFinish
				break
			}
			if err := sel.Process(); err != nil {
				vlog.ERROR("[%s] Process error %s", sel.name, err.Error())
			}
		}
	})
	return nil
}

func (sel *DataProcessorCtl) Process() error {
	// 获取最新区块
	latestNumber, err := sel.blockNum.GetBlockNumber()
	if err != nil {
		return err
	}
	// 1.获取交易数据
	dataList, err := sel.dataCursor.GetTxData()
	if err != nil {
		return err
	}
	vlog.INFO("获取到数量：%d", len(dataList))
	for _, data := range dataList {
		for _, prc := range sel.processorList {
			if err := prc.Process(int64(latestNumber), data); err != nil {
				// 记录错误的数据
				_ = sel.errorRepo.WriterErrorRecord("", data.TxHash, data.BlockNumber, err)
			}
		}
	}
	//
	_ = sel.dataCursor.Finish()
	return nil
}

func (sel *DataProcessorCtl) Exit(ctx context.Context) error {
	sel.isStop = true
	vlog.INFO("等待 [%s] 写入完成", sel.name)
	<-sel.waitFinish
	vlog.INFO("[%s] 写入完成", sel.name)
	close(sel.waitFinish)
	return nil
}
