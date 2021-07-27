package repo

import (
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/infra/dao"
)

//type BlockNumberRepo interface {
//	GetCntSyncBlockNumber() (int64, error)
//	UpdateSyncBlockNUmber(n int64) error
//}

func NewBlockNumberRepo(ebnDao *dao.EthereumBlockNumberDao) *BlockNumberRepo {
	return &BlockNumberRepo{ebnDao: ebnDao}
}

type BlockNumberRepo struct {
	ebnDao *dao.EthereumBlockNumberDao
}

// 初始最新区块
func (sel *BlockNumberRepo) InitLatestBlockNumber(bkNum int64) error {
	_, ok, err := sel.ebnDao.GetSyncBlockConfig("LatestBlockNumber")
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	return sel.ebnDao.SetSyncBlockConfig("LatestBlockNumber", fmt.Sprintf("%d", bkNum))
}

func (sel *BlockNumberRepo) UpdateLatestBlockNumber(bkNum int64) error {
	return sel.ebnDao.UpdateLatestBlockNumber(bkNum)
}

func (sel *BlockNumberRepo) GetCntSyncBlockNumber() (int64, error) {
	return sel.ebnDao.GetSyncBlockNumber()
}
func (sel *BlockNumberRepo) UpdateSyncBlockNUmber(n int64) error {
	return sel.ebnDao.UpdateSyncBlockNumber(n)
}
