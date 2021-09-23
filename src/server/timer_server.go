package server

import (
	"context"
)

type StartExit interface {
	Start() error
	Exit(ctx context.Context) error
}

type TimerServer struct {
	prc []StartExit
}

func NewTimerServer() *TimerServer {
	return &TimerServer{}
}

func (t *TimerServer) Add(s ...StartExit) {
	t.prc = append(t.prc, s...)
}

func (t *TimerServer) Schema() string {
	return "TimerServer"
}

func (t *TimerServer) Init(ctx context.Context) error {
	return nil
}

func (t *TimerServer) Start(ctx context.Context) error {
	for _, v := range t.prc {
		v.Start()
	}
	return nil
}

func (t *TimerServer) Exit(ctx context.Context) error {
	for _, v := range t.prc {
		v.Exit(ctx)
	}
	return nil
}
