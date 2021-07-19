package service

import (
	"context"
	"github.com/ville-vv/eth-chain-store/src/domain/ethm"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/vilgo/vlog"
	"time"
)

type SyncBlockNumberProgress interface {
	IsLatest(blockNumber int64) bool
	GetNow() (blockNumber int64, err error)
	Next()
}

// SyncBlockChainService 同步区块链数据服务
type SyncBlockChainService struct {
	ethMng       *ethm.EthereumManager
	syncInterval int64
	maxSyncNum   chan int
	syncPgs      SyncBlockNumberProgress
	stopCh       chan int
}

func NewSyncBlockChainService(ethRpcCli ethrpc.EthRPC, txWrite ethm.TxWriter, syncPgs SyncBlockNumberProgress) *SyncBlockChainService {
	s := &SyncBlockChainService{
		ethMng:       ethm.NewEthereumManager(ethRpcCli, txWrite),
		syncInterval: 0,
		maxSyncNum:   make(chan int, 100),
		syncPgs:      syncPgs,
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
	blockNumber, err := s.syncPgs.GetNow()
	if err != nil {
		vlog.ERROR("获取区块错误")
		return
	}
	if s.syncPgs.IsLatest(blockNumber) {
		if err = s.ethMng.PullBlock(); err != nil {
			vlog.ERROR("获取最新区块数据失败")
			return
		}
	} else {
		if err = s.ethMng.PullBlockByNumber(blockNumber); err != nil {
			vlog.ERROR("获取指定区块数据失败")
			return
		}
	}
	s.syncPgs.Next()
}
