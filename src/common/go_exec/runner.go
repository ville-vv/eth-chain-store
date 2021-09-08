package go_exec

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"
)

type Runner interface {
	Schema() string
	Init(ctx context.Context) error
	Start(ctx context.Context) error
	Exit(ctx context.Context) error
}

func Run(ctx context.Context, svr Runner) {
	if err := svr.Init(ctx); err != nil {
		return
	}
	Go(func() {
		if err := svr.Start(ctx); err != nil {
			os.Exit(1)
		}
	})

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-sig:
	case <-ctx.Done():
	}
	if err := svr.Exit(ctx); err != nil {
		return
	}
	time.Sleep(time.Second * 1)
	return
}

func Go(fn func()) {
	var gw sync.WaitGroup
	gw.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("%v \n %s", err, string(debug.Stack()))
			}
		}()
		gw.Done()
		fn()
	}()
	gw.Wait()
}
