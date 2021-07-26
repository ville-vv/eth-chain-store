package service

import (
	"context"
	"github.com/ville-vv/eth-chain-store/src/domain/ethm"
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
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

func NewSyncBlockChainService(ethRpcCli ethrpc.EthRPC, txWrite ethm.TxWriter, bkRepo repo.BlockNumberRepo) *SyncBlockChainService {
	syncCounter, err := ethm.NewSyncBlockNumberCounter(ethRpcCli, bkRepo)
	if err != nil {
		panic("NewSyncBlockChainService" + err.Error())
	}

	s := &SyncBlockChainService{
		ethMng:       ethm.NewEthereumManager(ethRpcCli, txWrite),
		syncInterval: 15,
		maxSyncNum:   make(chan int, 100),
		syncCounter:  syncCounter,
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
	s.fastSync()
	return nil
}

func (s *SyncBlockChainService) Exit(ctx context.Context) error {
	close(s.stopCh)
	return nil
}

// fastSync 快速的同步数据，间隔时间缩短
func (s *SyncBlockChainService) fastSync() {
	time.Sleep(time.Second)
	tk := time.NewTicker(time.Millisecond * 10)
	for {
		select {
		case <-tk.C:
			// 开启协程前添加一个控制，用于达到控制协程数量的目的
			if s.syncCounter.IsLatestBlockNumber() {
				goto startNormal
			}
			s.syncBlockChain()
		case <-s.stopCh:
			tk.Stop()
			return
		}
	}
startNormal:
	s.syncTimerTicker()
}

func (s *SyncBlockChainService) wait() {
	s.maxSyncNum <- 1
}

func (s *SyncBlockChainService) done() {
	<-s.maxSyncNum
}

func (s *SyncBlockChainService) syncTimerTicker() {
	// 区块打包是15秒中一次，所以 syncInterval 应该设置为 15 秒
	tk := time.NewTicker(time.Second * time.Duration(s.syncInterval))
	for {
		select {
		case <-tk.C:
			s.syncBlockChain()
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
	vlog.INFO("starting sync block [%d]", blockNumber)
	s.wait()
	go func(done func()) {
		defer done()
		// 执行完成后就释放一个
		if err = s.ethMng.PullBlockByNumber(blockNumber); err != nil {
			vlog.ERROR("获取指定区块数据失败 %d %s", blockNumber, err.Error())
			return
		}
		if err = s.syncCounter.FinishThisSync(blockNumber); err != nil {
			vlog.ERROR("更新同步的区块号失败 %s", err.Error())
		}
		vlog.INFO("finished sync block [%d]", blockNumber)
	}(s.done)
}
