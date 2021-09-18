package server

import (
	"errors"
	"github.com/ville-vv/vilgo/runner"
)

type ConfigOption struct {
}

type FactoryForServer struct {
}

func (sel *FactoryForServer) BuildServers(mould string) ([]runner.Runner, error) {
	switch mould {
	case "ToMysql":
	case "ToHiveFile":
	case "ToHive":
	}
	return nil, errors.New("not support this mould")
}
