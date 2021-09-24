package dao

import (
	"context"
	"github.com/ville-vv/vilgo/vlog"
	"time"
)

type TickTask struct {
	name      string
	exec      func()
	interval  time.Duration
	stopCh    chan int
	waithStop chan int
}

func NewTickTask(name string, interval time.Duration, exec func()) *TickTask {
	return &TickTask{name: name, exec: exec, interval: interval, stopCh: make(chan int), waithStop: make(chan int)}
}

func (t *TickTask) SetExec(exec func()) {
	t.exec = exec
}

func (t *TickTask) Scheme() string {
	return t.name
}

func (t *TickTask) Init() error {
	return nil
}

func (t *TickTask) Start() error {
	tmk := time.NewTicker(t.interval)
	for {
		select {
		case <-t.stopCh:
			t.exec()
			t.waithStop <- 1
			return nil
		case <-tmk.C:
			t.exec()
		}
	}
}

func (t *TickTask) Exit(ctx context.Context) error {
	close(t.stopCh)
	vlog.INFO("waiting tick task %s exist", t.name)
	<-t.waithStop
	return nil
}
