package service

import (
	"context"
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/common/conf"
	"github.com/ville-vv/eth-chain-store/src/domain/ethm"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/vilgo/vlog"
	"strconv"
	"sync"
	"time"
)

// SyncBlockChainService 同步区块链数据服务
type SyncBlockChainService struct {
	ethMng            *ethm.EthereumManager
	syncCounter       *ethm.SyncBlockNumberCounter
	syncInterval      time.Duration
	firstSyncInterval time.Duration
	maxSyncNum        chan int
	stopCh            chan int
	isStop            bool
}

func NewSyncBlockChainService(maxNum int, ethRpcCli ethrpc.EthRPC, txWrite ethm.TxWriter, bkRepo ethm.SyncBlockNumberPersist) *SyncBlockChainService {
	syncCounter, err := ethm.NewSyncBlockNumberCounter(ethRpcCli, bkRepo)
	if err != nil {
		panic("NewSyncBlockChainService" + err.Error())
	}

	var firstSyncIntervalStr = "1"
	var syncIntervalStr = "15"
	conf.ReadFlag(&firstSyncIntervalStr, "fsi")
	conf.ReadFlag(&syncIntervalStr, "fsi")

	firstSyncInterval, err := strconv.ParseInt(firstSyncIntervalStr, 10, 64)
	if err != nil {
		panic("NewSyncBlockChainService" + err.Error())
	}

	syncInterval, err := strconv.ParseInt(syncIntervalStr, 10, 64)
	if err != nil {
		panic("NewSyncBlockChainService" + err.Error())
	}

	s := &SyncBlockChainService{
		ethMng:            ethm.NewEthereumManager(ethRpcCli, txWrite),
		syncInterval:      time.Duration(syncInterval),
		maxSyncNum:        make(chan int, maxNum),
		syncCounter:       syncCounter,
		firstSyncInterval: time.Duration(firstSyncInterval),
		stopCh:            make(chan int),
	}

	return s
}

func (s *SyncBlockChainService) Scheme() string {
	return "SyncBlockChainService"
}

func (s *SyncBlockChainService) Init() error {
	// 获取当前已经同步的区块
	return nil
}

func (s *SyncBlockChainService) Start() error {
	s.fastSync()
	return nil
}

func (s *SyncBlockChainService) Exit(ctx context.Context) error {
	close(s.stopCh)
	s.isStop = true
	vlog.INFO("同步工具退出exit")
	return nil
}

// fastSync 快速的同步数据，间隔时间缩短
func (s *SyncBlockChainService) fastSync() {
	time.Sleep(time.Second)
	tk := time.NewTicker(time.Millisecond * s.firstSyncInterval)
	vlog.INFO("启动区块快速同步服务 internal is %d", s.firstSyncInterval)
	for {
		select {
		case <-tk.C:
			// 开启协程前添加一个控制，用于达到控制协程数量的目的
			if s.syncCounter.IsLatestBlockNumber() {
				goto startNormal
			}
			for {
				if s.isStop {
					vlog.INFO("同步工具退出")
					return
				}
				if s.syncBlockChainV2() {
					break
				}
				time.Sleep(time.Millisecond * 100)
			}
		case <-s.stopCh:
			vlog.INFO("同步工具退出stop")
			tk.Stop()
			return
		}
	}
startNormal:
	s.syncTimerTicker()
}

func (s *SyncBlockChainService) apply() {
	s.maxSyncNum <- 1
}

func (s *SyncBlockChainService) release() {
	<-s.maxSyncNum
}

func (s *SyncBlockChainService) syncTimerTicker() {
	// 区块打包是15秒中一次，所以 syncInterval 应该设置为 15 秒
	tk := time.NewTicker(time.Second * s.syncInterval)
	vlog.INFO("启动区块正常化同步服务 internal is %d", s.syncInterval)
	for {
		select {
		case <-tk.C:
			for {
				if s.isStop {
					return
				}
				if s.syncBlockChainV2() {
					break
				}
				time.Sleep(time.Millisecond * 100)
			}
		case <-s.stopCh:
			tk.Stop()
			return
		}
	}
}

// syncBlockChain 同步数据
func (s *SyncBlockChainService) syncBlockChain() {
	// 获取当前同步的区块
	blockNumber, err := s.syncCounter.GetSyncBlockNumber()
	if err != nil {
		vlog.ERROR("获取区块错误 %s", err.Error())
		return
	}
	latestBlockNumber := s.syncCounter.GetLatestBlockNumber()
	gw := sync.WaitGroup{}
	gw.Add(1)
	s.apply()
	go func(wg *sync.WaitGroup, done func()) {
		wg.Done()
		defer done()
		// 执行完成后就释放一个
		vlog.INFO("starting sync block [%d]", blockNumber)
		if err = s.ethMng.PullBlockByNumber(blockNumber, fmt.Sprintf("%d", latestBlockNumber)); err != nil {
			vlog.ERROR("获取指定区块数据失败 %d %s", blockNumber, err.Error())
			return
		}
		if err = s.syncCounter.FinishThisSync(blockNumber); err != nil {
			vlog.ERROR("更新同步的区块号失败 %s", err.Error())
		}
		vlog.INFO("finished sync block [%d]", blockNumber)
	}(&gw, s.release)
	gw.Wait()
}

// syncBlockChainV2 同步数据
func (s *SyncBlockChainService) syncBlockChainV2() bool {
	// 获取当前同步的区块
	blockNumber, err := s.syncCounter.GetSyncBlockNumber()
	if err != nil {
		vlog.ERROR("获取区块错误 %s", err.Error())
		return false
	}
	latestBlockNumber := s.syncCounter.GetLatestBlockNumber()
	if s.isStop {
		return false
	}
	// 执行完成后就释放一个
	vlog.INFO("starting sync block [%d]", blockNumber)
	if err = s.ethMng.PullBlockByNumber(blockNumber, fmt.Sprintf("%d", latestBlockNumber)); err != nil {
		vlog.ERROR("获取指定区块数据失败 %d %s", blockNumber, err.Error())
		return false
	}
	if err = s.syncCounter.FinishThisSync(blockNumber); err != nil {
		vlog.ERROR("更新同步的区块号失败 %s", err.Error())
	}
	vlog.INFO("finished sync block [%d]", blockNumber)
	return true
}
