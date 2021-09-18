package dao

import (
	"context"
	"time"
)

type TickTask struct {
	name     string
	exec     func()
	interval time.Duration
	stopCh   chan int
}

func NewTickTask(name string, interval time.Duration, exec func()) *TickTask {
	return &TickTask{name: name, exec: exec, interval: interval, stopCh: make(chan int)}
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
			return nil
		case <-tmk.C:
			t.exec()
		}
	}
}

func (t *TickTask) Exit(ctx context.Context) error {
	close(t.stopCh)
	return nil
}
