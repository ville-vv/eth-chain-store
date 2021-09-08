package ethm

import (
	"context"
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/common/go_exec"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/eth-chain-store/src/infra/monitor"
	"github.com/ville-vv/vilgo/vlog"
)

type TxWriterSchema interface {
	Schema() string
	TxWriter
}

// EthWriterControl 这里决定分发给哪些需要处理交易数据的 schema
// 控制最大可以开启多少个协程同时写
type EthWriterControl struct {
	goExec    *go_exec.GoExec
	txWriters map[string]TxWriterSchema
	isStop    bool
}

func NewEthWriterControl(maxNum int, txWriterList ...TxWriterSchema) *EthWriterControl {
	er := &EthWriterControl{
		goExec:    go_exec.New(maxNum),
		txWriters: make(map[string]TxWriterSchema),
		isStop:    false,
	}
	if len(txWriterList) > 0 {
		er.RegisterTxWriter(txWriterList...)
	}
	return er
}

func (sel *EthWriterControl) RegisterTxWriter(txwList ...TxWriterSchema) {
	for i := 0; i < len(txwList); i++ {
		txw := txwList[i]
		sel.txWriters[txw.Schema()] = txw
	}
}

func (sel *EthWriterControl) Scheme() string {
	return "EthWriterControl"
}

func (sel *EthWriterControl) Init() error {
	return nil
}

func (sel *EthWriterControl) Start() error {
	return nil
}

func (sel *EthWriterControl) Exit(ctx context.Context) error {
	sel.waitStop()
	return nil
}

func (sel *EthWriterControl) waitStop() {
	sel.goExec.WaitFinish(func(n int64) {
		vlog.INFO("eth data writer thread waiting close %d", n)
	})
}

func (sel *EthWriterControl) TxWrite(txData *model.TransactionData) error {
	if txData == nil {
		return nil
	}
	if sel.isStop {
		return fmt.Errorf("write process is stop")
	}
	sel.goExec.Go(txData, func(val interface{}) {
		monitor.TxWriteProcessNum.Inc()
		data := val.(*model.TransactionData)
		sel.txWrite(data)
		monitor.TxWriteProcessNum.Dec()
	})
	return nil
}

func (sel *EthWriterControl) txWrite(txData *model.TransactionData) {
	for _, txW := range sel.txWriters {
		txW.TxWrite(txData)
	}
}
