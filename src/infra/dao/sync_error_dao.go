package dao

import "github.com/ville-vv/eth-chain-store/src/infra/model"

type SyncErrorDao struct {
	db DB
}

func NewSyncErrorDao(db DB) *SyncErrorDao {
	return &SyncErrorDao{db: db}
}

func (sel *SyncErrorDao) WriterErrorRecord(info *model.SyncErrorRecord) error {
	return sel.db.GetDB().Create(info).Error
}
