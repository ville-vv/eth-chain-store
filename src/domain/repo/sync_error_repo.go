package repo

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/ville-vv/eth-chain-store/src/infra/dao"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"io"
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

type SyncErrorRepositoryV2 struct {
	w io.Writer
}

func NewSyncErrorRepositoryV2(fs io.Writer) *SyncErrorRepositoryV2 {
	return &SyncErrorRepositoryV2{w: fs}
}

func (sel *SyncErrorRepositoryV2) WriterErrorRecord(label string, bkNumber string, hash string, err error) error {
	data, err := jsoniter.Marshal(&model.SyncErrorRecord{
		Label:       label,
		BlockNumber: bkNumber,
		Hash:        hash,
		Err:         err.Error(),
	})
	_, err = sel.w.Write(data)
	return err
}
