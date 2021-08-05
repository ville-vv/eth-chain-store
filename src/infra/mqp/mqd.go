package mqp

import (
	"context"
	"fmt"
	"github.com/ville-vv/vilgo/runner"
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
}

func NewMDP(logf LogFunc) *MQD {
	if logf == nil {
		logf = logPrintf
	}
	m := &MQD{
		msgChan:   make(chan *Message, 10000),
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
	return nil
}

func (sel *MQD) Exit(ctx context.Context) error {
	close(sel.exitChan)
	fmt.Println("MQD Exit")
	//sel.RLock()
	//for _, c := range sel.csmMap {
	//	retrySvc, ok := c.(*retryConsumer)
	//	if ok {
	//		_ = retrySvc.Exit(context.Background())
	//	}
	//}
	//sel.RUnlock()
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
	return
}

func (sel *MQD) Publish(msg *Message) error {
	select {
	case sel.msgChan <- msg:
	}
	return nil
}

func (sel *MQD) SubScribe(cum Consumer) {
	sel.Lock()
	sel.csmMap[cum.ID()] = cum
	sel.Unlock()
	return
}
