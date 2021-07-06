package eth_native

import "time"

type TbEthAddress struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;COMMENT:" name:""`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;COMMENT:" name:""`
	EthAddress
}

type EthAddress struct {
}
