package eth_native

import "time"

// 合约地址表
type TbContractAddress struct {
	ID        int64     `gorm:"primary_key"`
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
