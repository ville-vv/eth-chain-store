package ethm

import (
	"context"
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/vilgo/vlog"
	"github.com/ville-vv/vilgo/vtask"
)

type reTryElem struct {
	cntRetry int
	data     interface{}
}

func newReTryElem(val interface{}) vtask.RetryElem {
	return &reTryElem{
		cntRetry: 1,
		data:     val,
	}
}

func (r *reTryElem) Can() bool {
	if r.cntRetry <= 0 {
		return false
	}
	r.cntRetry--
	return true
}

func (r *reTryElem) Interval() int64 {
	return 1
}

func (r *reTryElem) GetData() interface{} {
	return r.data
}

type RetryProcessor struct {
	mTask     *vtask.MiniTask
	errorRepo repo.SyncErrorRepository
	txWriter  TxWriter
	name      string
	isStop    bool
	stopCh    chan int
}

func NewRetryProcess(name string, txWriter TxWriter, errorRepo repo.SyncErrorRepository) *RetryProcessor {
	er := &RetryProcessor{
		name: name, txWriter: txWriter,
		errorRepo: errorRepo,
		stopCh:    make(chan int),
	}
	er.mTask = vtask.NewMiniTask(&vtask.TaskOption{
		RetryFlag:       true,
		NewRetry:        newReTryElem,
		ErrEventHandler: er.ErrWriter,
		Exec:            er.exec,
	})
	return er
}

func (sel *RetryProcessor) ErrWriter(ctx interface{}, err error) {
	data, ok := ctx.(*model.TransactionData)
	if !ok {
		vlog.ERROR("ctx:%v, error:%s", ctx, err.Error())
		return
	}
	if Err := sel.errorRepo.WriterErrorRecord("", data.BlockNumber, data.Hash, err); Err != nil {
		vlog.ERROR("writer error information hash:%s error %s", data.Hash, Err.Error())
	}
}

func (sel *RetryProcessor) exec(val interface{}) (retry bool) {
	data, ok := val.(*model.TransactionData)
	if !ok {
		return false
	}
	//vlog.WARN("retry to write tx data %s", data.Hash)
	if err := sel.txWriter.TxWrite(data); err != nil {
		//vlog.WARN("retry to write tx data %s %s", data.Hash, err.Error())
		return true
	}
	return false
}

func (sel *RetryProcessor) Scheme() string {
	return sel.name
}

func (sel *RetryProcessor) Schema() string {
	return sel.name
}

func (sel *RetryProcessor) Init() error {
	return nil
}

func (sel *RetryProcessor) Exit(ctx context.Context) error {
	sel.mTask.Stop()
	close(sel.stopCh)
	vlog.INFO("retry processor exited [%s]", sel.name)
	return nil
}

func (sel *RetryProcessor) Start() error {
	sel.mTask.Start()
	vlog.INFO("retry processor started [%s]", sel.name)
	//for i := 0; i < sel.maxNum; i++ {
	//	go func() {
	//		for {
	//			select {
	//			case data, ok := <-sel.dataChan:
	//				if !ok {
	//					return
	//				}
	//				_ = sel.TxWrite(data)
	//				sel.runCounter.Dec()
	//			}
	//		}
	//	}()
	//}
	return nil
}

func (sel *RetryProcessor) TxWrite(txData *model.TransactionData) error {
	err := sel.txWriter.TxWrite(txData)
	if err != nil {
		vlog.WARN("tx write failed push to retry %s", err.Error())
		if err = sel.mTask.Push(txData); err != nil {
			vlog.ERROR("eth retry push error %s", err.Error())
		}
	}
	return nil
}
