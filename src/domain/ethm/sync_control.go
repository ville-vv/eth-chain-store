package ethm

import (
	"context"
	"github.com/ville-vv/eth-chain-store/src/common/conf"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/vilgo/runner"
	"github.com/ville-vv/vilgo/vlog"
	"github.com/ville-vv/vilgo/vtask"
	"sync"
	"time"
)

type SyncConfigDataPersist interface {
	// InitLatestBlockNumber
	InitLatestBlockNumber(bkNum int64) error

	UpdateLatestBlockNumber(bkNum int64) error
	// GetCntSyncBlockNumber 获取当前同步了最区块号
	GetCntSyncBlockNumber() (int64, error)
	// UpdateSyncBlockNUmber 更新当前同步了的区块号
	UpdateSyncBlockNUmber(n int64) error
}

type SyncBlockPullFunc func(cntBKNum int64, laterBKNum int64) (int64, error)

func (sel SyncBlockPullFunc) Pull(cntBKNum int64, laterBKNum int64) (int64, error) {
	return sel(cntBKNum, laterBKNum)
}

type SyncBlockPuller interface {
	Pull(cntBKNum int64, laterBKNum int64) (int64, error)
}

type OptionConfig struct {
	StartBlockNumber int64
	EndBlockNumber   int64
	MaxSyncThreads   int
	EthRpcCli        ethrpc.EthRPC
	BknRepo          SyncConfigDataPersist
	Puller           SyncBlockPuller
}

type SyncBlockControl struct {
	lbnLock              sync.RWMutex
	ethRpcCli            ethrpc.EthRPC
	bknRepo              SyncConfigDataPersist
	latestNum            int64      // 最新区块
	intervalTotalNum     int64      // 间隔时间内总条数
	beforeSyncNumber     int64      // 上一次同步的区块号
	cntSyncBlockNumber   int64      // 当前同步的区块号
	startSyncBlockNumber int64      // 开始同步区块号
	endSyncBlockNumber   int64      // 同步结束的区块号， 如果为0就一直同步到最新区块号
	maxSyncThreads       int        // 最大同时同步数(每秒)
	syncBlockNumberChan  chan int64 // 同步的区块号分发
	//syncBlockNumberChanNoLimit chan int64
	stopCh     chan int
	syncStopCh chan int
	isStop     bool
	finishLock sync.Mutex
	threadNum  vtask.AtomicInt64
	puller     SyncBlockPuller
	rpcNotOpen bool
}

func (sel *SyncBlockControl) Scheme() string {
	return "SyncBlockControl"
}

func (sel *SyncBlockControl) Init() error {
	return nil
}

func NewSyncBlockControl(maxSyncNum int, ethRpcCli ethrpc.EthRPC, bknRepo SyncConfigDataPersist) *SyncBlockControl {
	sbn := &SyncBlockControl{
		ethRpcCli:           ethRpcCli,
		bknRepo:             bknRepo,
		maxSyncThreads:      maxSyncNum,
		syncBlockNumberChan: make(chan int64, maxSyncNum),
		stopCh:              make(chan int),
		syncStopCh:          make(chan int, maxSyncNum),
	}
	return sbn
}

func NewSyncBlockControlWithOpt(opt *OptionConfig) *SyncBlockControl {
	sbn := &SyncBlockControl{
		ethRpcCli:            opt.EthRpcCli,
		bknRepo:              opt.BknRepo,
		maxSyncThreads:       opt.MaxSyncThreads,
		endSyncBlockNumber:   opt.EndBlockNumber,
		startSyncBlockNumber: opt.StartBlockNumber,
		syncBlockNumberChan:  make(chan int64, opt.MaxSyncThreads),
		stopCh:               make(chan int),
		syncStopCh:           make(chan int, opt.MaxSyncThreads),
		puller:               opt.Puller,
	}
	return sbn
}

func (sel *SyncBlockControl) Start() error {
	cntNumber, err := sel.bknRepo.GetCntSyncBlockNumber()
	// 获取已经同步到的最近区块
	if sel.startSyncBlockNumber > cntNumber  {
		cntNumber = sel.startSyncBlockNumber
	}
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
		sel.loopGenSyncNumberOneSecond()
	})
	return nil
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
		sel.latestNum = int64(latestNumber)
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
	var diffNum = int64(0)
	var interval = int64(5)
	var intervalTotalNum int64
	for {
		select {
		case <-conf.GlobalExitSignal:
			return
		default:
		}
		time.Sleep(time.Second * time.Duration(interval))
		if !sel.rpcNotOpen {
			continue
		}

		sel.lbnLock.Lock()
		latestNumber = sel.latestNum
		sel.lbnLock.Unlock()

		if sel.latestNum > 0 {
			// 更新数据库
			if err := sel.bknRepo.UpdateLatestBlockNumber(sel.latestNum); err != nil {
				vlog.ERROR("get block number is error %s", err.Error())
			}
		}
		sel.finishLock.Lock()
		diffNum = sel.beforeSyncNumber - cntBKNum
		cntBKNum = sel.beforeSyncNumber
		intervalTotalNum = sel.intervalTotalNum
		sel.finishLock.Unlock()
		if cntBKNum > 0 {
			if cntBKNum >= sel.endSyncBlockNumber && sel.endSyncBlockNumber > 0 {
				break
			}
			vlog.INFO("blocks=[%d], record total num=[%d], current sync block no=[%d], latest=[%d]", diffNum, intervalTotalNum, cntBKNum, latestNumber)
			_ = sel.bknRepo.UpdateSyncBlockNUmber(sel.beforeSyncNumber)
		}
	}
	vlog.INFO("同步到最后区块，请主动结束")
}

// loopGenSyncNumber 每秒中生成区块数
func (sel *SyncBlockControl) loopGenSyncNumberOneSecond() {
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

// loopGenSyncNumberNoLimit 获取区块的速度取决于，写数据的数据， 这里不控制每秒取多少个块的数据
// 只要 syncBlockNumberChan 存在多余的缓存就直接获取
// TODO 未完成后续有时间再写
func (sel *SyncBlockControl) loopGenSyncNumberNoLimit() {
	for {
		if sel.endSyncBlockNumber > 0 && sel.endSyncBlockNumber < sel.cntSyncBlockNumber {
			// 同步结束了
			return
		}
		// 如果当前区块
		if sel.latestNum < sel.cntSyncBlockNumber {
			// 如果当前要同步的区块号大于最新的区块号就值跳出
			return
		}
		if sel.isStop {
			return
		}
		if !sel.rpcNotOpen {
			vlog.WARN("eth node rpc server is stop !!!!!!!!")
			time.Sleep(time.Second)
			return
		}
		sel.syncBlockNumberChan <- sel.cntSyncBlockNumber
		sel.cntSyncBlockNumber++
	}
}

func (sel *SyncBlockControl) genSyncBlockNumber() {
	// 一段时间内生成 maxSyncNun 需要同步的区块
	for index := 0; index < sel.maxSyncThreads; index++ {
		if sel.endSyncBlockNumber > 0 && sel.endSyncBlockNumber < sel.cntSyncBlockNumber {
			return
		}

		if sel.latestNum < sel.cntSyncBlockNumber {
			// 如果当前要同步的区块号大于最新的区块号就值跳出
			return
		}
		if sel.isStop {
			return
		}
		if !sel.rpcNotOpen {
			vlog.WARN("eth node rpc server is stop !!!!!!!!")
			time.Sleep(time.Second)
			return
		}

		sel.syncBlockNumberChan <- sel.cntSyncBlockNumber
		sel.cntSyncBlockNumber++
	}
}

func (sel *SyncBlockControl) SetPuller(pull SyncBlockPuller) {
	sel.puller = pull
}

func (sel *SyncBlockControl) startPull() {
	var waitGp sync.WaitGroup
	waitGp.Add(sel.maxSyncThreads)
	for i := 0; i < sel.maxSyncThreads; i++ {
		go func(seq int) {
			var err error
			var bkNumber int64
			var totalNum int64
			waitGp.Done()
			sel.threadNum.Inc()
			vlog.INFO("pull thread [%d] start", seq)
			for {
				select {
				// 定时生成需要同步的区块号，在这里捕捉
				case bkNumber = <-sel.syncBlockNumberChan:
					//vlog.INFO("同步区块：%d", bkNumber)
					totalNum, err = sel.puller.Pull(bkNumber, sel.GetLatestBlockNumber())
					if err != nil {
						vlog.ERROR("sync block number data failed %d %s", bkNumber, err.Error())
						continue
					}
					if err = sel.FinishSync(bkNumber, totalNum); err != nil {
						vlog.ERROR("update finished sync block number error %s", err.Error())
						continue
					}
					//vlog.INFO("同步区块完成：%d", bkNumber)
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

func (sel *SyncBlockControl) SetSyncFunc(ctx context.Context, syncFunc func(cntBKNum, laterBKNum int64) (int64, error)) {
	var waitGp sync.WaitGroup
	waitGp.Add(sel.maxSyncThreads)
	for i := 0; i < sel.maxSyncThreads; i++ {
		go func(seq int) {
			var err error
			var bkNumber int64
			var totalNum int64
			waitGp.Done()
			sel.threadNum.Inc()
			vlog.INFO("pull thread [%d] start", seq)
			for {
				select {
				// 定时生成需要同步的区块号，在这里捕捉
				case bkNumber = <-sel.syncBlockNumberChan:
					totalNum, err = syncFunc(bkNumber, sel.GetLatestBlockNumber())
					if err != nil {
						vlog.ERROR("sync block number data failed %d %s", bkNumber, err.Error())
						continue
					}
					if err = sel.FinishSync(bkNumber, totalNum); err != nil {
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
	lbn := sel.latestNum
	sel.lbnLock.RUnlock()
	return lbn
}

func (sel *SyncBlockControl) FinishSync(bkNum int64, totalNum int64) error {
	sel.finishLock.Lock()
	defer sel.finishLock.Unlock()
	if bkNum > sel.beforeSyncNumber {
		sel.beforeSyncNumber = bkNum
	}
	sel.intervalTotalNum += totalNum
	return nil
}

func (sel *SyncBlockControl) WaitExit() {

	for {
		if sel.threadNum.Load() <= 0 {
			return
		}
	}
}

func (sel *SyncBlockControl) Exit(ctx context.Context) error {
	if sel.isStop {
		return nil
	}
	sel.isStop = true
	close(sel.stopCh)
	close(sel.syncStopCh)
	if sel.latestNum > 0 {
		// 更新数据库
		if err := sel.bknRepo.UpdateLatestBlockNumber(sel.latestNum); err != nil {
			vlog.ERROR("get block number is error %s", err.Error())
		}
	}
	if sel.beforeSyncNumber > 0 {
		_ = sel.bknRepo.UpdateSyncBlockNUmber(sel.beforeSyncNumber)
	}

	sel.WaitExit()
	vlog.INFO("sync block number control exited")
	return nil
}
