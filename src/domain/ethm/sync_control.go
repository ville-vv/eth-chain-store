package ethm

import (
	"context"
	"github.com/ville-vv/eth-chain-store/src/common/conf"
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/vilgo/runner"
	"github.com/ville-vv/vilgo/vlog"
	"github.com/ville-vv/vilgo/vtask"
	"sync"
	"time"
)

type SyncBlockNumberPersist interface {
	//
	InitLatestBlockNumber(bkNum int64) error

	UpdateLatestBlockNumber(bkNum int64) error
	// 获取当前同步了最区块号
	GetCntSyncBlockNumber() (int64, error)
	// 更新当前同步了的区块号
	UpdateSyncBlockNUmber(n int64) error
}

// 区块同步计数器
type SyncBlockNumberCounter struct {
	haveDoneLock     sync.Mutex
	syncLock         sync.RWMutex
	finishLock       sync.Mutex
	cntSyncingNumber int64
	haveDoneCap      int
	haveDoneIndex    int
	haveDoneList     []int64
	latestNumber     int64
	ethRpcCli        ethrpc.EthRPC
	persist          SyncBlockNumberPersist
	beforeSyncNumber int64
}

func NewSyncBlockNumberCounter(ethRpcCli ethrpc.EthRPC, persist SyncBlockNumberPersist) (*SyncBlockNumberCounter, error) {
	cntNumber, err := persist.GetCntSyncBlockNumber()
	//cntNumber = 12696310 // 12696310
	if err != nil {
		return nil, err
	}
	s := &SyncBlockNumberCounter{
		cntSyncingNumber: cntNumber + 1,
		haveDoneCap:      1000,
		haveDoneList:     make([]int64, 1000),
		ethRpcCli:        ethRpcCli,
		persist:          persist,
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
			err = sel.persist.InitLatestBlockNumber(sel.latestNumber)
			if err != nil {
				panic(err)
			}
		}
		// 更新数据库
		if err = sel.persist.UpdateLatestBlockNumber(sel.latestNumber); err != nil {
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

func (sel *SyncBlockNumberCounter) GetLatestBlockNumber() int64 {
	sel.syncLock.RLock()
	latestNumber := sel.latestNumber
	sel.syncLock.RUnlock()
	return latestNumber
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
	sel.finishLock.Lock()
	defer sel.finishLock.Unlock()
	can := blockNumber > sel.beforeSyncNumber
	if can {
		sel.beforeSyncNumber = blockNumber
	}
	if can {
		return sel.persist.UpdateSyncBlockNUmber(sel.beforeSyncNumber)
	}
	return nil
}

type SyncBlockPullFunc func(cntBKNum int64, laterBKNum int64) error

func (sel SyncBlockPullFunc) Pull(cntBKNum int64, laterBKNum int64) error {
	return sel(cntBKNum, laterBKNum)
}

type SyncBlockPuller interface {
	Pull(cntBKNum int64, laterBKNum int64) error
}

type SyncBlockControl struct {
	lbnLock   sync.RWMutex
	ethRpcCli ethrpc.EthRPC
	//syncFinishList      []int64
	bknRepo             *repo.BlockNumberRepo
	lbNum               int64      // 最新区块
	cntSyncBlockNumber  int64      // 当前同步的区块号
	maxSyncNum          int        // 最大同时同步数(每秒)
	syncBlockNumberChan chan int64 // 同步的区块号分发
	stopCh              chan int
	syncStopCh          chan int
	isStop              bool
	finishLock          sync.Mutex
	beforeSyncNumber    int64
	threadNum           vtask.AtomicInt64
	puller              SyncBlockPuller

	rpcNotOpen bool
}

func NewSyncBlockControl(maxSyncNum int, ethRpcCli ethrpc.EthRPC, bknRepo *repo.BlockNumberRepo) *SyncBlockControl {
	sbn := &SyncBlockControl{
		ethRpcCli:           ethRpcCli,
		bknRepo:             bknRepo,
		maxSyncNum:          maxSyncNum,
		syncBlockNumberChan: make(chan int64, maxSyncNum),
		stopCh:              make(chan int),
		syncStopCh:          make(chan int, maxSyncNum),
	}
	return sbn
}

func (sel *SyncBlockControl) Start() {
	// 获取已经同步到的最近区块
	cntNumber, err := sel.bknRepo.GetCntSyncBlockNumber()
	//cntNumber = 12696310 // 12696310
	if err != nil {
		panic("NewSyncBlockControl " + err.Error())
	}
	sel.cntSyncBlockNumber = cntNumber

	var waitChan = make(chan int)
	go sel.loopSyncLatestBlockNumber(waitChan)
	<-waitChan

	go sel.loopSyncConfigUpdate()

	sel.startPull()

	runner.Go(func() {
		sel.loopGenSyncNumber()
	})
}

// loopSyncLatestBlockNumber 获取最新区块号
func (sel *SyncBlockControl) loopSyncLatestBlockNumber(waitStart chan int) {
	first := true
	oldNumber := uint64(0)
	for {
		time.Sleep(time.Second)
		if sel.isStop {
			return
		}
		latestNumber, err := sel.ethRpcCli.GetBlockNumber()
		if err != nil {
			sel.rpcNotOpen = false
			vlog.ERROR("get latest block number from main chain %s", err.Error())
			continue
		}
		sel.rpcNotOpen = true
		if oldNumber >= latestNumber {
			continue
		}

		sel.lbnLock.Lock()
		sel.lbNum = int64(latestNumber)
		sel.lbnLock.Unlock()
		if first {
			first = false
			waitStart <- 1
			err = sel.bknRepo.InitLatestBlockNumber(int64(latestNumber))
			if err != nil {
				panic(err)
			}
			// 更新数据库
			if err = sel.bknRepo.UpdateLatestBlockNumber(int64(latestNumber)); err != nil {
				vlog.ERROR("get block number is error %s", err.Error())
			}
		}

		oldNumber = latestNumber
	}
}

func (sel *SyncBlockControl) loopSyncConfigUpdate() {
	var cntBKNum = int64(0)
	var latestNumber = int64(0)

	for {
		select {
		case <-conf.GlobalExitSignal:
			return
		default:
		}
		time.Sleep(time.Second * 5)
		if !sel.rpcNotOpen {
			continue
		}

		sel.lbnLock.Lock()
		latestNumber = sel.lbNum
		sel.lbnLock.Unlock()

		if sel.lbNum > 0 {
			// 更新数据库
			if err := sel.bknRepo.UpdateLatestBlockNumber(sel.lbNum); err != nil {
				vlog.ERROR("get block number is error %s", err.Error())
			}
		}
		sel.finishLock.Lock()
		cntBKNum = sel.beforeSyncNumber
		sel.finishLock.Unlock()
		if cntBKNum > 0 {
			vlog.INFO("synchronized block number [%d], latest block number [%d]", cntBKNum, latestNumber)
			_ = sel.bknRepo.UpdateSyncBlockNUmber(sel.beforeSyncNumber)
		}
	}
}

// loopGenSyncNumber 没秒中生成区块数
func (sel *SyncBlockControl) loopGenSyncNumber() {
	tmr := time.NewTicker(time.Second)
	for {
		select {
		case <-tmr.C:
			if sel.isStop {
				return
			}
			sel.genSyncBlockNumber()
		case <-sel.stopCh:
			return
		}
	}
}

func (sel *SyncBlockControl) genSyncBlockNumber() {
	// 一段时间内生成maxSyncNun需要同步的区块
	for index := 0; index < sel.maxSyncNum; index++ {
		if sel.lbNum < sel.cntSyncBlockNumber {
			return
		}
		if sel.isStop {
			return
		}
		if !sel.rpcNotOpen {
			vlog.WARN("eth node rpc server is stop !!!!!!!!")
			return
		}

		sel.syncBlockNumberChan <- sel.cntSyncBlockNumber
		sel.cntSyncBlockNumber++
	}
}

//// 订阅同步区块
//func (sel *SyncBlockControl) SubSyncNumberChan() <-chan int64 {
//	return sel.syncBlockNumberChan
//}

func (sel *SyncBlockControl) SetPuller(pull SyncBlockPuller) {
	sel.puller = pull
}

func (sel *SyncBlockControl) startPull() {
	var waitGp sync.WaitGroup
	waitGp.Add(sel.maxSyncNum)
	for i := 0; i < sel.maxSyncNum; i++ {
		go func(seq int) {
			var err error
			var bkNumber int64
			waitGp.Done()
			sel.threadNum.Inc()
			vlog.INFO("pull thread [%d] start", seq)
			for {
				select {
				// 定时生成需要同步的区块号，在这里捕捉
				case bkNumber = <-sel.syncBlockNumberChan:
					err = sel.puller.Pull(bkNumber, sel.GetLatestBlockNumber())
					if err != nil {
						vlog.ERROR("sync block number data failed %d %s", bkNumber, err.Error())
						continue
					}
					if err = sel.FinishSync(bkNumber); err != nil {
						vlog.ERROR("update finished sync block number error %s", err.Error())
						continue
					}
				case <-sel.syncStopCh:
					goto exit
				}
			}
		exit:
			sel.threadNum.Dec()
			vlog.INFO("pull thread [%d] exited", seq)
		}(i)
	}
	waitGp.Wait()
}

func (sel *SyncBlockControl) SetSyncFunc(ctx context.Context, syncFunc func(cntBKNum, laterBKNum int64) error) {
	var waitGp sync.WaitGroup
	waitGp.Add(sel.maxSyncNum)
	for i := 0; i < sel.maxSyncNum; i++ {
		go func(seq int) {
			var err error
			var bkNumber int64
			waitGp.Done()
			sel.threadNum.Inc()
			vlog.INFO("pull thread [%d] start", seq)
			for {
				select {
				// 定时生成需要同步的区块号，在这里捕捉
				case bkNumber = <-sel.syncBlockNumberChan:
					err = syncFunc(bkNumber, sel.GetLatestBlockNumber())
					if err != nil {
						vlog.ERROR("sync block number data failed %d %s", bkNumber, err.Error())
						continue
					}
					if err = sel.FinishSync(bkNumber); err != nil {
						vlog.ERROR("update finished sync block number error %s", err.Error())
						continue
					}
				case <-sel.syncStopCh:
					goto exit
				case <-ctx.Done():
					goto exit
				}
			}
		exit:
			sel.threadNum.Dec()
			vlog.INFO("pull thread [%d] exit", seq)
		}(i)
	}
	waitGp.Wait()
}

func (sel *SyncBlockControl) GetLatestBlockNumber() int64 {
	sel.lbnLock.RLock()
	lbn := sel.lbNum
	sel.lbnLock.RUnlock()
	return lbn
}

func (sel *SyncBlockControl) FinishSync(bkNum int64) error {
	sel.finishLock.Lock()
	defer sel.finishLock.Unlock()
	if bkNum > sel.beforeSyncNumber {
		sel.beforeSyncNumber = bkNum
	}
	return nil
}

func (sel *SyncBlockControl) WaitExit() {

	for {
		if sel.threadNum.Load() <= 0 {
			return
		}
	}
}

func (sel *SyncBlockControl) Exit() {
	sel.isStop = true
	close(sel.stopCh)
	close(sel.syncStopCh)
	// 更新数据库
	if err := sel.bknRepo.UpdateLatestBlockNumber(sel.lbNum); err != nil {
		vlog.ERROR("get block number is error %s", err.Error())
	}
	_ = sel.bknRepo.UpdateSyncBlockNUmber(sel.beforeSyncNumber)
	sel.WaitExit()
}
