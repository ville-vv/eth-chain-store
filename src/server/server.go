package server

import (
	"context"
	"github.com/ville-vv/vilgo/runner"
)

type Server struct {
	runner []runner.Runner
}

func (s *Server) Schema() string {
	return "Server"
}

func (s *Server) Init(ctx context.Context) error {
	for _, r := range s.runner {
		_ = r.Init()
	}
	return nil
}

func (s *Server) Start(ctx context.Context) error {
	for _, r := range s.runner {
		runner.Go(func() {
			_ = r.Start()
		})
	}
	return nil
}

func (s *Server) Exit(ctx context.Context) error {
	for _, r := range s.runner {
		_ = r.Exit(ctx)
	}
	return nil
}

func (s *Server) Add(r ...runner.Runner) {
	s.runner = append(s.runner, r...)
}
