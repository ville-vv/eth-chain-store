package ethm

import (
	"context"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/vilgo/vlog"
	"github.com/ville-vv/vilgo/vtask"
)

type reTryElem struct {
	cntRetry int
	data     interface{}
}

func (r *reTryElem) Can() bool {
	if r.cntRetry < 0 {
		return false
	}
	r.cntRetry--
	return true
}

func (r *reTryElem) Interval() int64 {
	return 5
}

func (r *reTryElem) GetData() interface{} {
	return r.data
}

type EthRetryWriter struct {
	mTask    *vtask.MiniTask
	txWriter TxWriter
}

func NewEthRetryWriter(txWriter TxWriter) *EthRetryWriter {
	mTask := vtask.NewMiniTask(&vtask.TaskOption{
		RetryFlag: true,
		NewRetry: func(val interface{}) vtask.RetryElem {
			return &reTryElem{
				cntRetry: 4,
				data:     val,
			}
		},
		ErrEventHandler: func(ctx interface{}, err error) {
			vlog.ERROR("ctx:%v, error:%s", ctx, err.Error())
		},
		Exec: func(val interface{}) (retry bool) {
			data, ok := val.(*model.TransactionData)
			if !ok {
				return false
			}
			err := txWriter.TxWrite(data)
			if err != nil {
				return true
			}
			return false
		},
	})
	return &EthRetryWriter{txWriter: txWriter, mTask: mTask}
}

func (sel *EthRetryWriter) Scheme() string {
	return "EthRetryWriter"
}

func (sel *EthRetryWriter) Init() error {
	return nil
}

func (sel *EthRetryWriter) Exit(ctx context.Context) error {
	sel.mTask.Stop()
	return nil
}

func (sel *EthRetryWriter) Start() error {
	sel.mTask.Start()
	return nil
}

func (sel *EthRetryWriter) TxWrite(txData *model.TransactionData) error {
	err := sel.txWriter.TxWrite(txData)
	if err != nil {
		if err = sel.mTask.Push(txData); err != nil {
			vlog.ERROR("eth retry push error %s", err.Error())
		}
	}
	return nil
}
