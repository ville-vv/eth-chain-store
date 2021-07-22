package dao

import (
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"strconv"
)

type EthereumBlockNumberDao struct {
	db DB
}

func NewEthereumBlockNumberDao(db DB) *EthereumBlockNumberDao {
	return &EthereumBlockNumberDao{db: db}
}

func (sel *EthereumBlockNumberDao) GetSyncBlockNumber() (int64, error) {
	var syncInfo model.SyncBlockConfig
	err := sel.db.GetDB().Debug().Model(&syncInfo).Select("value").Where("k_name='SyncBlockNumber'").First(&syncInfo).Error
	if err != nil {
		return 0, err
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
