package mqp

import (
	"context"
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

type retryConsumer struct {
	mTask *vtask.MiniTask
	Consumer
}

func newRetryConsumer(logf LogFunc, csm Consumer) *retryConsumer {
	rm := &retryConsumer{Consumer: csm}
	mTask := vtask.NewMiniTask(&vtask.TaskOption{
		RetryFlag: true,
		NewRetry: func(val interface{}) vtask.RetryElem {
			return &reTryElem{
				cntRetry: 4,
				data:     val,
			}
		},
		ErrEventHandler: func(ctx interface{}, err error) {
			logf("%v %s", ctx, ctx)
		},
		Exec: rm.exec,
	})
	mTask.Start()
	rm.mTask = mTask
	return rm
}

func (sel *retryConsumer) exec(val interface{}) (retry bool) {
	data, ok := val.(*Message)
	if !ok {
		return false
	}
	err := sel.Consumer.Process(data)
	if err != nil {
		return true
	}
	return false
}

func (sel *retryConsumer) Process(msg *Message) error {
	err := sel.Consumer.Process(msg)
	if err != nil {
		return sel.mTask.Push(msg)
	}
	return nil
}

func (sel *retryConsumer) Exit(ctx context.Context) error {
	sel.mTask.Stop()
	return nil
}
