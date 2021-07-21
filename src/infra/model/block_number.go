package model

import "time"

// TbBlockNumber 区块信息
type TbBlockNumber struct {
	ID          uint        `gorm:"primary_key"`
	CreatedAt   time.Time   `json:"created_at" gorm:"column:created_at;COMMENT:" name:""`
	UpdatedAt   time.Time   `json:"updated_at" gorm:"column:updated_at;COMMENT:" name:""`
	BlockNumber BlockNumber `json:"block_number" gorm:"column:block_number;COMMENT:" name:""`
}

type BlockNumber uint64

type SyncBlockConfig struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;COMMENT:" name:""`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;COMMENT:" name:""`
	KName     string    `json:"k_name" gorm:"column:k_name;type:varchar(255);COMMENT:" name:""`
	Value     string    `json:"value" gorm:"column:value;type:varchar(255);COMMENT:" name:""`
}
