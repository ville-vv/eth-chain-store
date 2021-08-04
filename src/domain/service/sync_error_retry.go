package service

import (
	"context"
	"github.com/ville-vv/vilgo/vlog"
	"github.com/ville-vv/vilgo/vtask"
)

type SyncErrorRetryService struct {
	mTask *vtask.MiniTask
}

func (s *SyncErrorRetryService) Scheme() string {
	return "SyncErrorRetryService"
}

func (s *SyncErrorRetryService) Init() error {
	s.mTask = vtask.NewMiniTask(&vtask.TaskOption{
		RetryFlag:  true,
		Persistent: nil,
		NewRetry:   nil,
		ErrEventHandler: func(ctx interface{}, err error) {
			vlog.ERROR("ctx:%v, error:%s", ctx, err.Error())
		},
		Exec: nil,
	})
	return nil
}

func (s *SyncErrorRetryService) Start() error {
	s.mTask.Start()
	return nil
}

func (s *SyncErrorRetryService) Exit(ctx context.Context) error {
	s.mTask.Stop()
	return nil
}
