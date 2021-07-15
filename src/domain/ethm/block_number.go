package ethm

import (
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/vilgo/vlog"
)

type BlockNumberCaptor struct {
	cntBlockNumber uint64
	ethRequester   ethrpc.EthRPC
	bkRepo         repo.IBlockNumberRepo
}

// BlockNumberUpdater 从链上获取最新区块号并更新
func (p *BlockNumberCaptor) UpdateBlockNumber() {
	// 拉取数据
	bkNumber, err := p.ethRequester.GetBlockNumber()
	if err != nil {
		vlog.ERROR("get block number is error %s", err.Error())
		return
	}
	p.cntBlockNumber = bkNumber
	// 更新数据库
	if err = p.bkRepo.UpdateBlockNumber(bkNumber); err != nil {
		vlog.ERROR("get block number is error %s", err.Error())
		return
	}
}

// GetBlockNumber 获取当前区块号
func (p *BlockNumberCaptor) GetBlockNumber() uint64 {
	return p.cntBlockNumber
}
