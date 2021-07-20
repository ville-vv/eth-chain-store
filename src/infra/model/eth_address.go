package model

import "time"

// 原生以太坊地址信息
type EthereumAccount struct {
	ID          int64     `gorm:"primary_key"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;COMMENT:" name:""`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;COMMENT:" name:""`
	Address     string    `json:"address" gorm:"column:address;COMMENT:" name:""`
	FirstTxTime string    `json:"first_transfer_time" gorm:"column:first_transfer_time;COMMENT:" name:""`
	FirstTxHash string    `json:"first_tx_hash" gorm:"column:first_tx_hash;COMMENT:" name:""`
	Balance     string    `json:"balance" gorm:"column:balance;COMMENT:" name:""`
}

// 合约地址表
type TbContractAddress struct {
	ID        int64     `gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;COMMENT:" name:""`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;COMMENT:" name:""`
	ContractContent
}

// 合约地址信息
type ContractContent struct {
	Symbol      string `json:"symbol" gorm:"column:symbol;COMMENT:" name:""`             // 合约代号
	Address     string `json:"address" gorm:"column:address;COMMENT:" name:""`           // 合约地址
	PublishTime string `json:"publish_time" gorm:"column:publish_time;COMMENT:" name:""` // 合约发布时间
	//EthBalance  string `json:"balance" gorm:"column:balance;COMMENT:" name:""`           // 合约Eth余额
	//TokenBalance string `json:"contract_balance" gorm:"column:contract_balance;COMMENT:" name:""` // 合约代币余额
	IsErc20     bool   `json:"is_erc_20" gorm:"column:is_erc_20;COMMENT:" name:""`       // 是否为ERC20代币
	TotalSupply string `json:"total_supply" gorm:"column:total_supply;COMMENT:" name:""` // 代币总发行量
}

// 合约账户
type AccountContent struct {
	Contract string `json:"contract" name:""`                               // 账户的合约地址
	Address  string `json:"address" gorm:"column:address;COMMENT:" name:""` // 账户地址
	TxTime   string `json:"tx_time" name:""`                                // 交易时间
	TxHash   string `json:"tx_hash" name:""`                                // 交易 Hash
	Balance  string `json:"balance" gorm:"column:balance;COMMENT:" name:""` // 账户余额
	//IsContract   bool   `json:"is_contract" name:""`                                              // 是否为合约代币
	Symbol       string `json:"symbol" gorm:"column:symbol;COMMENT:" name:""`                     // 合约代号
	TokenBalance string `json:"contract_balance" gorm:"column:contract_balance;COMMENT:" name:""` // 合约代币余额
	IsErc20      bool   `json:"is_erc_20" gorm:"column:is_erc_20;COMMENT:" name:""`               // 是否为ERC20代币
	//TotalSupply     string `json:"total_supply" gorm:"column:total_supply;COMMENT:" name:""`         // 代币总发行量
}
