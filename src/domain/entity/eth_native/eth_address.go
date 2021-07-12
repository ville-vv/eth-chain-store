package eth_native

import "time"

// 所有以太坊地址表
type TbEthAddress struct {
	ID        int64     `gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;COMMENT:" name:""`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;COMMENT:" name:""`
	NativeAddressInfo
}

// 原生以太坊地址信息
type NativeAddressInfo struct {
	Address     string `json:"address" gorm:"column:address;COMMENT:" name:""`
	FirstTxTime string `json:"first_transfer_time" gorm:"column:first_transfer_time;COMMENT:" name:""`
	FirstTxHash string `json:"first_tx_hash" gorm:"column:first_tx_hash;COMMENT:" name:""`
	Balance     string `json:"balance" gorm:"column:balance;COMMENT:" name:""`
}
