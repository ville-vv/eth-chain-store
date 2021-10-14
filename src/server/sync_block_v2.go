package server

import (
	"context"
	"github.com/ville-vv/eth-chain-store/src/domain/ethm"
	"github.com/ville-vv/vilgo/vlog"
)

type SyncBlockChainServiceV2 struct {
	stopCh  chan int
	ethMng  *ethm.EthereumDataPuller
	counter *ethm.SyncBlockControl
}

func NewSyncBlockChainServiceV2(ethMng *ethm.EthereumDataPuller, counter *ethm.SyncBlockControl) *SyncBlockChainServiceV2 {
	sbc := &SyncBlockChainServiceV2{ethMng: ethMng, counter: counter}
	return sbc
}

func (sel *SyncBlockChainServiceV2) Scheme() string {
	return "SyncBlockChainServiceV2"
}

func (sel *SyncBlockChainServiceV2) Init() error {
	return nil
}

func (sel *SyncBlockChainServiceV2) Start() error {
	sel.counter.SetPuller(sel)
	sel.counter.Start()
	return nil
}

func (sel *SyncBlockChainServiceV2) Exit(ctx context.Context) error {
	sel.counter.Exit(ctx)
	return nil
}

func (sel *SyncBlockChainServiceV2) Pull(cntBKNum int64, laterBKNum int64) (int64, error) {
	n, err := sel.ethMng.Pull(cntBKNum, laterBKNum)
	if err != nil {
		vlog.ERROR("sync block number data failed %d %s", cntBKNum, err.Error())
		return n, err
	}
	return n, nil
}
