package async

import (
	"context"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/vilgo/vlog"
	"sync"
	"time"
)

type LatestBlockNumberCache struct {
	sync.RWMutex
	latestBlockNumber uint64
	rpcCli            ethrpc.EthRPC
	isStop            bool
}

func NewLatestBlockNumberCache(rpcCli ethrpc.EthRPC) *LatestBlockNumberCache {
	return &LatestBlockNumberCache{rpcCli: rpcCli}
}

func (sel *LatestBlockNumberCache) Start() error {
	var waitChan = make(chan int)
	var isFirst = true
	go func() {
		for {
			if sel.isStop {
				break
			}
			bkn, err := sel.rpcCli.GetBlockNumber()
			if err != nil {
				vlog.ERROR("获取最新区块错误....")
				time.Sleep(time.Second * 3)
				continue
			}
			sel.Lock()
			sel.latestBlockNumber = bkn
			sel.Unlock()
			if isFirst {
				waitChan <- 1
			}

		}
	}()
	vlog.INFO("等待获取最新区块中....")
	<-waitChan
	return nil
}

func (sel *LatestBlockNumberCache) Exit(ctx context.Context) error {
	sel.isStop = true
	return nil
}

func (sel *LatestBlockNumberCache) GetBlockNumber() (uint64, error) {
	return sel.latestBlockNumber, nil
}
