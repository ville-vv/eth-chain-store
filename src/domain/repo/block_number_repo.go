package repo

import "github.com/ville-vv/eth-chain-store/src/infra/dao"

type BlockNumberRepo interface {
	UpdateBlockNumber(bkNum int64) error
	GetCntSyncBlockNumber() (int64, error)
	UpdateSyncBlockNUmber(n int64) error
}

func NewBlockNumberRepo(ebnDao *dao.EthereumBlockNumberDao) BlockNumberRepo {
	return &blockNumberRepo{ebnDao: ebnDao}
}

type blockNumberRepo struct {
	ebnDao *dao.EthereumBlockNumberDao
}

func (sel *blockNumberRepo) UpdateBlockNumber(bkNum int64) error {
	return nil
}
func (sel *blockNumberRepo) GetCntSyncBlockNumber() (int64, error) {
	return 0, nil
}
func (sel *blockNumberRepo) UpdateSyncBlockNUmber(n int64) error {
	return nil
}
