package model

import "time"

type SyncErrorRecord struct {
	ID          int64     `gorm:"primary_key"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;COMMENT:" name:""`
	Label       string    `json:"label" gorm:"column:label;COMMENT:" name:""`
	BlockNumber string    `json:"block_number" gorm:"column:block_number;COMMENT:" name:""`
	Hash        string    `json:"hash" gorm:"column:hash;COMMENT:" name:""`
	Err         string    `json:"err" gorm:"column:err;type:varchar(255);COMMENT:" name:""`
}
