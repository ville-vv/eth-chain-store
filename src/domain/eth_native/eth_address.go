package eth_native

import "time"

// 所有以太坊地址表
type TbEthAddress struct {
	ID        uint      `gorm:"primary_key"`
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

// 合约地址表
type TbContractAddress struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;COMMENT:" name:""`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;COMMENT:" name:""`
	ContractAddressInfo
}

// 合约地址信息
type ContractAddressInfo struct {
	Symbol       string `json:"symbol" gorm:"column:symbol;COMMENT:" name:""`                     // 合约代号
	Address      string `json:"address" gorm:"column:address;COMMENT:" name:""`                   // 合约地址
	PublishTime  string `json:"publish_time" gorm:"column:publish_time;COMMENT:" name:""`         // 合约发布时间
	EthBalance   string `json:"balance" gorm:"column:balance;COMMENT:" name:""`                   // 合约Eth余额
	TokenBalance string `json:"contract_balance" gorm:"column:contract_balance;COMMENT:" name:""` // 合约代币余额
	IsErc20      bool   `json:"is_erc_20" gorm:"column:is_erc_20;COMMENT:" name:""`               // 是否为ERC20代币
	TotalSupply  string `json:"total_supply" gorm:"column:total_supply;COMMENT:" name:""`         // 代币总发行量
}

// 合约代币的配置表，不然不知道
type TbContractConfig struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;COMMENT:" name:""`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;COMMENT:" name:""`
	Address   string    `json:"address" gorm:"column:address;COMMENT:" name:""` // 合约地址
	Symbol    string    `json:"symbol" gorm:"column:symbol;COMMENT:" name:""`   //
	Issuer    string    `json:"issuer" gorm:"column:issuer;COMMENT:" name:""`
}

// 地址合约代币绑定表，如果一个地址有过合约交易就会在这表里面生成一个绑定关系
type TbContractBindAddress struct {
	ID              uint      `gorm:"primary_key"`
	CreatedAt       time.Time `json:"created_at" gorm:"column:created_at;COMMENT:" name:""`
	UpdatedAt       time.Time `json:"updated_at" gorm:"column:updated_at;COMMENT:" name:""`
	Address         string    `json:"address" gorm:"column:address;COMMENT:" name:""`
	ContractAddress string    `json:"contract_address" gorm:"column:contract_address;COMMENT:" name:""`
	Symbol          string    `json:"symbol" gorm:"column:symbol;COMMENT:" name:""`
	Balance         string    `json:"balance" gorm:"column:balance;COMMENT:" name:""`
}
