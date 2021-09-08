package model

import "time"

// 原生以太坊地址信息
type EthereumAccount struct {
	ID          int64     `gorm:"primary_key"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;COMMENT:" name:""`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;COMMENT:" name:""`
	Address     string    `json:"address" gorm:"column:address;index;type:varchar(255);COMMENT:" name:""`
	FirstTxTime string    `json:"first_tx_time" gorm:"column:first_tx_time;COMMENT:" name:""`
	FirstTxHash string    `json:"first_tx_hash" gorm:"column:first_tx_hash;COMMENT:" name:""`
	Balance     string    `json:"balance" gorm:"column:balance;type:varchar(255);COMMENT:" name:""`
}

// 合约地址表
type ContractAddressRecord struct {
	ID        int64     `gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;COMMENT:" name:""`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;COMMENT:" name:""`
	ContractContent
}

// 合约地址信息
type ContractContent struct {
	Symbol      string `json:"symbol" gorm:"column:symbol;type:varchar(255) CHARACTER SET 'utf8mb4' COLLATE 'utf8mb4_general_ci';COMMENT:" name:""` // 合约代号
	Address     string `json:"address" gorm:"column:address;index;type:varchar(127);COMMENT:" name:""`                                              // 合约地址
	PublishTime string `json:"publish_time" gorm:"column:publish_time;type:varchar(255);COMMENT:" name:""`                                          // 合约发布时间
	//EthBalance  string `json:"balance" gorm:"column:balance;COMMENT:" name:""`           // 合约Eth余额
	//TokenBalance string `json:"contract_balance" gorm:"column:contract_balance;COMMENT:" name:""` // 合约代币余额
	IsErc20     bool   `json:"is_erc20" gorm:"column:is_erc20;COMMENT:" name:""`                           // 是否为ERC20代币
	TotalSupply string `json:"total_supply" gorm:"column:total_supply;type:varchar(255);COMMENT:" name:""` // 代币总发行量
	DecimalBit  int    `json:"decimal_bit" gorm:"column:decimal_bit;COMMENT:" name:""`                     // 小数位
}
