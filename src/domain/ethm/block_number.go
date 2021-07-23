package ethm

import (
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/vilgo/vlog"
	"sync"
	"time"
)

type SyncBlockNumberCounter struct {
	haveDoneLock     sync.Mutex
	syncLock         sync.Mutex
	cntSyncingNumber int64
	haveDoneCap      int
	haveDoneIndex    int
	haveDoneList     []int64
	latestNumber     int64
	ethRpcCli        ethrpc.EthRPC
	bkRepo           repo.BlockNumberRepo
	beforeSyncNumber int64
}

func NewSyncBlockNumberCounter(ethRpcCli ethrpc.EthRPC, bkRepo repo.BlockNumberRepo) (*SyncBlockNumberCounter, error) {
	cntNumber, err := bkRepo.GetCntSyncBlockNumber()
	cntNumber = 12696216 // 12696216
	if err != nil {
		return nil, err
	}
	s := &SyncBlockNumberCounter{
		haveDoneLock:     sync.Mutex{},
		syncLock:         sync.Mutex{},
		cntSyncingNumber: cntNumber,
		haveDoneCap:      100,
		haveDoneList:     make([]int64, 100),
		ethRpcCli:        ethRpcCli,
		bkRepo:           bkRepo,
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go s.loopSyncBlockNumber(&wg)
	wg.Wait()
	return s, nil
}

func (sel *SyncBlockNumberCounter) loopSyncBlockNumber(wg *sync.WaitGroup) {
	first := true
	for {
		latestNumber, err := sel.ethRpcCli.GetBlockNumber()
		if err != nil {
			vlog.ERROR("get latest block number from main chain %s", err.Error())
			time.Sleep(time.Second)
			continue
		}
		sel.syncLock.Lock()
		sel.latestNumber = int64(latestNumber)
		sel.syncLock.Unlock()
		if first {
			first = false
			wg.Done()
		}

		// 更新数据库
		if err = sel.bkRepo.UpdateBlockNumber(sel.latestNumber); err != nil {
			vlog.ERROR("get block number is error %s", err.Error())
		}
		time.Sleep(time.Second * 15)
	}
}

// IsLatestBlockNumber 是不是最新区块
func (sel *SyncBlockNumberCounter) IsLatestBlockNumber() bool {
	sel.syncLock.Lock()
	latestNumber := sel.latestNumber
	sel.syncLock.Unlock()
	return sel.cntSyncingNumber >= latestNumber
}

func (sel *SyncBlockNumberCounter) SetSyncing(blockNumber int64) bool {
	// 设置正在同步的区块
	sel.haveDoneLock.Lock()
	sel.haveDoneList[sel.haveDoneIndex] = blockNumber
	sel.haveDoneIndex++
	if sel.haveDoneIndex >= sel.haveDoneCap {
		sel.haveDoneIndex++
	}
	sel.haveDoneLock.Unlock()
	return false
}

func (sel *SyncBlockNumberCounter) IsSyncing(blockNumber int64) bool {
	for i := 0; i < sel.haveDoneCap; i++ {
		// 判断该区块是否正在同步
		if sel.haveDoneList[i] == blockNumber {
			return true
		}
	}
	return false
}

func (sel *SyncBlockNumberCounter) GetSyncBlockNumber() (blockNumber int64, err error) {
	sel.syncLock.Lock()
	blockNumber = sel.cntSyncingNumber
	sel.haveDoneList[sel.haveDoneIndex] = blockNumber
	sel.cntSyncingNumber++
	if sel.cntSyncingNumber > sel.latestNumber {
		sel.cntSyncingNumber = sel.latestNumber
	}
	sel.syncLock.Unlock()

	return blockNumber, nil
}

func (sel *SyncBlockNumberCounter) FinishThisSync(blockNumber int64) error {
	sel.syncLock.Lock()
	can := blockNumber > sel.beforeSyncNumber
	if can {
		sel.beforeSyncNumber = blockNumber
	}
	sel.syncLock.Unlock()
	if can {
		return sel.bkRepo.UpdateSyncBlockNUmber(sel.beforeSyncNumber)
	}
	return nil
}
