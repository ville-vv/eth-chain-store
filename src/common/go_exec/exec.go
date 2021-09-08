package go_exec

import (
	"fmt"
	"go.uber.org/atomic"
	"time"
)

type GoExec struct {
	maxThread      int
	maxThreadLimit chan int
	isFinish       bool
	threadCounter  atomic.Int64 //当前线程数
}

func New(maxThread int) *GoExec {
	return &GoExec{
		maxThread:      maxThread,
		maxThreadLimit: make(chan int, maxThread),
		threadCounter:  atomic.Int64{},
	}
}

func (sel *GoExec) Go(val interface{}, f func(interface{})) {
	if sel.isFinish {
		fmt.Println("go exec have exit")
		return
	}
	sel.maxThreadLimit <- 1
	sel.threadCounter.Inc()
	go func(dt interface{}) {
		f(dt)
		sel.threadCounter.Dec()
		<-sel.maxThreadLimit
	}(val)
}

func (sel *GoExec) Threads() int64 {
	return sel.threadCounter.Load()
}

func (sel *GoExec) WaitFinish(f func(n int64)) {
	num := int64(0)
	for {
		num = sel.threadCounter.Load()
		f(num)
		if num <= 0 {
			break
		}
		time.Sleep(time.Millisecond * 10)
	}
	sel.isFinish = true
	close(sel.maxThreadLimit)
}
