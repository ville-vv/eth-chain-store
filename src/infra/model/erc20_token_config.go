package model

import "time"

// 合约代币的配置表，不然不知道 jas
type Erc20TokenConfig struct {
	ID        int64     `gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;COMMENT:" name:""`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;COMMENT:" name:""` // 更新时间
	Erc20TokenConfigContent
}

type Erc20TokenConfigContent struct {
	Address string `json:"address" gorm:"column:address;COMMENT:" name:""` // 合约地址
	Symbol  string `json:"symbol" gorm:"column:symbol;COMMENT:" name:""`   // 合约代币标识
	Issuer  string `json:"issuer" gorm:"column:issuer;COMMENT:" name:""`   // 发行者
}
