package mqp

import (
	"context"
	"errors"
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/infra/monitor"
	"github.com/ville-vv/vilgo/runner"
	"github.com/ville-vv/vilgo/vtask"
	"net/http"
	"sync"
)

type LogFunc func(format string, args ...interface{})

type MQD struct {
	sync.RWMutex
	msgChan   chan *Message
	updateCSM chan int
	exitChan  chan int
	csmMap    map[string]Consumer
	logf      LogFunc
	isStop    bool
	poolSize  vtask.AtomicInt64
}

func NewMDP(pooSize int, logf LogFunc) *MQD {
	if logf == nil {
		logf = logPrintf
	}
	m := &MQD{
		msgChan:   make(chan *Message, pooSize),
		updateCSM: make(chan int),
		exitChan:  make(chan int),
		csmMap:    make(map[string]Consumer),
		logf:      logf,
	}
	return m
}

func (sel *MQD) Scheme() string {
	return "MQD"
}

func (sel *MQD) Init() error {
	return nil
}

func (sel *MQD) Start() error {
	runner.Go(sel.msgPump)
	//go sel.msgPump()
	//runner.Go(sel.monitor)
	return nil
}

func (sel *MQD) monitor() {
	http.HandleFunc("/mq_size", func(writer http.ResponseWriter, request *http.Request) {
		bodyStr := fmt.Sprintf(`
Pool Size: %d
`, sel.poolSize.Load())
		writer.Write([]byte(bodyStr))
	})
	http.ListenAndServe("0.0.0.0:5489", nil)
}

func (sel *MQD) Exit(ctx context.Context) error {
	sel.isStop = true
	close(sel.exitChan)
	close(sel.updateCSM)
	sel.ClearMsgChan()
	sel.logf("MQD Exited")
	return nil
}

func logPrintf(format string, args ...interface{}) {
	fmt.Printf(format, args)
}

func (sel *MQD) msgPump() {
	var (
		msg     *Message
		msgChan chan *Message
		csms    []Consumer
	)
	sel.RLock()
	for _, c := range sel.csmMap {
		csms = append(csms, c)
	}
	sel.RUnlock()
	if len(sel.csmMap) > 0 {
		msgChan = sel.msgChan
	}

	for {
		select {
		case msg = <-msgChan:
			monitor.MqSize.Dec()
		case <-sel.updateCSM:
			csms = csms[:0]
			sel.RLock()
			for _, c := range sel.csmMap {
				csms = append(csms, c)
			}
			sel.RUnlock()
			if len(csms) == 0 {
				msgChan = nil
			}
			if len(sel.csmMap) > 0 {
				msgChan = sel.msgChan
			}
		case <-sel.exitChan:
			goto exit
		}
		for i, c := range csms {
			tempMsg := msg
			if i > 0 {
				tempMsg = msg.Copy()
			}
			err := c.Process(tempMsg)
			if err != nil {
				sel.logf("[MQD] consumer:%s error:%s", c.ID(), err.Error())
			}
		}
	}
exit:
	sel.logf("MQD Exiting")
	return
}

func (sel *MQD) Publish(msg *Message) error {
	if sel.isStop {
		return errors.New("mqd have closed")
	}
	select {
	case sel.msgChan <- msg:
		monitor.MqSize.Inc()
	}
	return nil
}

func (sel *MQD) ClearMsgChan() {
	close(sel.msgChan)
	sel.logf("clear msg ")
	for v := range sel.msgChan {
		monitor.MqSize.Dec()
		for _, csm := range sel.csmMap {
			err := csm.Process(v)
			if err != nil {
				sel.logf("[MQD] consumer:%s error:%s", csm.ID(), err.Error())
			}
		}
	}
}

func (sel *MQD) SubScribe(cum Consumer) {
	sel.Lock()
	sel.csmMap[cum.ID()] = cum
	sel.Unlock()
	return
}
