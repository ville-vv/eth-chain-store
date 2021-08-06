package repo

import (
	"github.com/ville-vv/eth-chain-store/src/infra/dao"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
)

type SyncErrorRepository interface {
	WriterErrorRecord(label string, bkNumber string, hash string, err error) error
}

func NewSyncErrorRepository(errorDao *dao.SyncErrorDao) SyncErrorRepository {
	return &syncErrorRepository{errorDao: errorDao}
}

type syncErrorRepository struct {
	errorDao *dao.SyncErrorDao
}

func (sel *syncErrorRepository) WriterErrorRecord(label string, bkNumber string, hash string, err error) error {
	return sel.errorDao.WriterErrorRecord(&model.SyncErrorRecord{
		Label:       label,
		BlockNumber: bkNumber,
		Hash:        hash,
		Err:         err.Error(),
	})
}
