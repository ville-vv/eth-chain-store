package dao

import (
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"gorm.io/gorm"
	"strconv"
)

type EthereumBlockNumberDao struct {
	db DB
}

func NewEthereumBlockNumberDao(db DB) *EthereumBlockNumberDao {
	return &EthereumBlockNumberDao{db: db}
}

func (sel *EthereumBlockNumberDao) SetSyncBlockConfig(k, v string) error {
	return sel.db.GetDB().Create(&model.SyncBlockConfig{KName: k, Value: v}).Error
}

func (sel *EthereumBlockNumberDao) GetSyncBlockConfig(k string) (string, bool, error) {
	sbc := &model.SyncBlockConfig{}
	err := sel.db.GetDB().Select("value").Where("k_name=?", k).First(sbc).Error
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			return sbc.Value, false, err
		}
		// 没有记录
		return "", false, nil
	}
	return sbc.Value, true, nil
}

func (sel *EthereumBlockNumberDao) UpdateLatestBlockNumber(blockNumber int64) error {
	db := sel.db.GetDB()
	return db.Model(&model.SyncBlockConfig{}).Where("k_name='LatestBlockNumber'").Update("value", blockNumber).Error
}

func (sel *EthereumBlockNumberDao) GetSyncBlockNumber() (int64, error) {
	var syncInfo model.SyncBlockConfig
	err := sel.db.GetDB().Model(&syncInfo).Select("value").Where("k_name='SyncBlockNumber'").First(&syncInfo).Error
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			return 0, err
		}
		// 没有记录先创建
		syncInfo.Value = "0"
		syncInfo.KName = "SyncBlockNumber"
		return 0, sel.db.GetDB().Create(&syncInfo).Error
	}
	blockNumber, err := strconv.ParseInt(syncInfo.Value, 0, 64)
	if err != nil {
		return 0, err
	}
	return blockNumber, nil
}

func (sel *EthereumBlockNumberDao) UpdateSyncBlockNumber(blockNumber int64) error {
	db := sel.db.GetDB()
	return db.Model(&model.SyncBlockConfig{}).Where("k_name='SyncBlockNumber'").Update("value", blockNumber).Error
}
