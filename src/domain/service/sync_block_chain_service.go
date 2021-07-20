package service

import (
	"context"
	"github.com/ville-vv/eth-chain-store/src/domain/ethm"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/vilgo/vlog"
	"time"
)

// SyncBlockChainService 同步区块链数据服务
type SyncBlockChainService struct {
	ethMng       *ethm.EthereumManager
	syncCounter  *ethm.SyncBlockNumberCounter
	syncInterval int64
	maxSyncNum   chan int
	stopCh       chan int
}

func NewSyncBlockChainService(ethRpcCli ethrpc.EthRPC, txWrite ethm.TxWriter, syncPgs *ethm.SyncBlockNumberCounter) *SyncBlockChainService {
	s := &SyncBlockChainService{
		ethMng:       ethm.NewEthereumManager(ethRpcCli, txWrite),
		syncInterval: 0,
		maxSyncNum:   make(chan int, 100),
		syncCounter:  syncPgs,
		stopCh:       make(chan int),
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
	s.syncTimerTicker()
	return nil
}

func (s *SyncBlockChainService) Exit(ctx context.Context) error {
	close(s.stopCh)
	return nil
}

// fastSync 快速的同步数据，间隔时间缩短
func (s *SyncBlockChainService) fastSync() {
	tk := time.NewTicker(time.Second)
	for {
		select {
		case <-tk.C:
			// 开启协程前添加一个控制，用于达到控制协程数量的目的
			if s.syncCounter.IsLatestBlockNumber() {
				goto startNormal
			}
			s.maxSyncNum <- 1
			go func() {
				s.syncBlockChain()
				// 执行完成后就释放一个
				<-s.maxSyncNum
			}()
		case <-s.stopCh:
			tk.Stop()
			return
		}
	}
startNormal:
	s.syncTimerTicker()
}

func (s *SyncBlockChainService) syncTimerTicker() {
	// 区块打包是15秒中一次，所以 syncInterval 应该设置为 15 秒
	tk := time.NewTicker(time.Second * time.Duration(s.syncInterval))
	for {
		select {
		case <-tk.C:
			// 开启协程前添加一个控制，用于达到控制协程数量的目的
			s.maxSyncNum <- 1
			go func() {
				s.syncBlockChain()
				// 执行完成后就释放一个
				<-s.maxSyncNum
			}()
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
		vlog.ERROR("获取区块错误")
		return
	}
	if s.syncCounter.IsSyncing(blockNumber) {
		return
	}
	if err = s.ethMng.PullBlockByNumber(blockNumber); err != nil {
		vlog.ERROR("获取指定区块数据失败")
		return
	}
	s.syncCounter.FinishThisSync(blockNumber)
}
