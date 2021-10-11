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
	trigger   chan int
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
			vlog.INFO("waiting task %s existed", t.name)
			t.exec()
			vlog.INFO("tick task %s existed", t.name)
			close(t.waithStop)
			return nil
		case <-tmk.C:
			t.exec()
		case <-t.trigger:
			t.exec()
		}
	}
}

func (t *TickTask) Exit(ctx context.Context) error {
	vlog.INFO("waiting tick task %s exist, after 5 minute not any response will force exited", t.name)
	close(t.stopCh)
	select {
	case <-time.After(time.Minute * 5):
		return nil
	case <-t.waithStop:
		return nil
	}
}

func (t *TickTask) Trigger() {
	t.trigger <- 1
}
