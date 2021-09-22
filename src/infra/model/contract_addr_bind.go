package model

import "time"

// 地址合约代币绑定表，如果一个地址有过合约交易就会在这表里面生成一个绑定关系
// 如果数据库表数据过大，就另外分表
type ContractAccountBind struct {
	ID              int64      `gorm:"primary_key"`
	CreatedAt       time.Time  `json:"created_at" gorm:"column:created_at;COMMENT:" name:""`                                                          // 记录创建时间
	UpdatedAt       time.Time  `json:"updated_at" gorm:"column:updated_at;COMMENT:" name:""`                                                          // 记录更新时间
	TxTime          *time.Time `json:"tx_time" gorm:"column:tx_time;COMMENT:" name:""`                                                                // 第一笔交易时间
	Address         string     `json:"address" gorm:"column:address;type:varchar(128);index:idx_cab_addr_crtaddr;COMMENT:" name:""`                   // 钱包地址
	ContractAddress string     `json:"contract_address" gorm:"column:contract_address;type:varchar(128);index:idx_cab_addr_crtaddr;COMMENT:" name:""` // 合约代币地址
	Symbol          string     `json:"symbol" gorm:"column:symbol;type:varchar(128);COMMENT:" name:""`                                                // 合约代币标识
	Balance         string     `json:"balance" gorm:"column:balance;type:varchar(128);COMMENT:" name:""`                                              // 合约代币金额
}
