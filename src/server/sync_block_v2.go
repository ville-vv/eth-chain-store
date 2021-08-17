package server

import (
	"context"
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/domain/ethm"
	"github.com/ville-vv/vilgo/vlog"
)

type SyncBlockChainServiceV2 struct {
	stopCh  chan int
	ethMng  *ethm.EthereumManager
	counter *ethm.SyncBlockControl
}

func NewSyncBlockChainServiceV2(ethMng *ethm.EthereumManager, counter *ethm.SyncBlockControl) *SyncBlockChainServiceV2 {
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
	//sel.counter.SetSyncFunc(context.Background(), sel.doneSync)
	return nil
}

func (sel *SyncBlockChainServiceV2) Exit(ctx context.Context) error {
	sel.counter.Exit()
	vlog.INFO("同步工具退出exit")
	return nil
}

func (sel *SyncBlockChainServiceV2) doneSync(bkNumber, latestBlockNumber int64) error {
	var err error
	//vlog.INFO("starting sync block [%d]", bkNumber)
	if err = sel.ethMng.PullBlockByNumber(bkNumber, fmt.Sprintf("%d", latestBlockNumber)); err != nil {
		vlog.ERROR("sync block number data failed %d %s", bkNumber, err.Error())
		return err
	}
	vlog.INFO("finished sync block [%d]", bkNumber)
	return nil
}

func (sel *SyncBlockChainServiceV2) Pull(cntBKNum int64, laterBKNum int64) error {
	var err error
	//vlog.INFO("starting sync block [%d]", bkNumber)
	if err = sel.ethMng.PullBlockByNumber(cntBKNum, fmt.Sprintf("%d", laterBKNum)); err != nil {
		vlog.ERROR("sync block number data failed %d %s", cntBKNum, err.Error())
		return err
	}
	return nil
}
