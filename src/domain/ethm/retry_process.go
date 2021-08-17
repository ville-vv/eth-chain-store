package ethm

import (
	"context"
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/eth-chain-store/src/infra/mqp"
	"github.com/ville-vv/vilgo/vlog"
	"github.com/ville-vv/vilgo/vtask"
	"time"
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

type RetryProcess struct {
	mTask      *vtask.MiniTask
	errorRepo  repo.SyncErrorRepository
	txWriter   TxWriter
	name       string
	maxGoRun   chan int
	runCounter *vtask.AtomicInt64
	isStop     bool
}

func NewRetryProcess(name string, maxNum int, txWriter TxWriter, errorRepo repo.SyncErrorRepository) *RetryProcess {
	er := &RetryProcess{name: name, txWriter: txWriter, maxGoRun: make(chan int, maxNum), errorRepo: errorRepo}
	er.mTask = vtask.NewMiniTask(&vtask.TaskOption{
		RetryFlag:       true,
		NewRetry:        newReTryElem,
		ErrEventHandler: er.ErrWriter,
		Exec:            er.exec,
	})
	return er
}

func (sel *RetryProcess) SetMonitor(runCounter *vtask.AtomicInt64) {
	sel.runCounter = runCounter
}

func (sel *RetryProcess) ErrWriter(ctx interface{}, err error) {
	data, ok := ctx.(*model.TransactionData)
	if !ok {
		vlog.ERROR("ctx:%v, error:%s", ctx, err.Error())
		return
	}
	if Err := sel.errorRepo.WriterErrorRecord("", data.BlockNumber, data.Hash, err); Err != nil {
		vlog.ERROR("writer error information hash:%s error %s", data.Hash, Err.Error())
	}
}

func (sel *RetryProcess) exec(val interface{}) (retry bool) {
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

func (sel *RetryProcess) Scheme() string {
	return "RetryProcess"
}

func (sel *RetryProcess) Init() error {
	return nil
}

func (sel *RetryProcess) Exit(ctx context.Context) error {
	sel.mTask.Stop()
	sel.waitStop()
	return nil
}

func (sel *RetryProcess) waitStop() {
	for {
		vlog.INFO("Wait Close RetryProcess %s %d", sel.name, sel.runCounter.Load())
		if sel.runCounter.Load() <= 0 {
			vlog.INFO("Exit RetryProcess %s ", sel.name)
			return
		}

		time.Sleep(time.Second)
	}
}

func (sel *RetryProcess) Start() error {
	sel.mTask.Start()
	return nil
}

func (sel *RetryProcess) TxWrite(txData *model.TransactionData) error {
	err := sel.txWriter.TxWrite(txData)
	if err != nil {
		vlog.WARN("tx write failed push to retry %s", err.Error())
		if err = sel.mTask.Push(txData); err != nil {
			vlog.ERROR("eth retry push error %s", err.Error())
		}
	}
	return nil
}

func (sel *RetryProcess) ID() string {
	return sel.name
}

func (sel *RetryProcess) Process(msg *mqp.Message) error {
	if msg == nil {
		return nil
	}
	if sel.isStop {
		return fmt.Errorf("%s write process is stop", sel.name)
	}
	txData := &model.TransactionData{}
	err := msg.UnMarshalFromBody(txData)
	if err != nil {
		return err
	}

	sel.maxGoRun <- 1
	sel.runCounter.Inc()
	go func(dt *model.TransactionData) {
		_ = sel.TxWrite(dt)
		<-sel.maxGoRun
		sel.runCounter.Dec()
	}(txData)
	return nil
}
