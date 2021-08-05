package model

import "time"

const (
	TxEventTypeTransfer = "transfer"
)

const (
	TxTypeNormal    = "normal transactions"
	TxTypeInternal  = "internal transactions"
	TxTokenTransfer = "token transfers"
)

type TransactionData struct {
	LatestNumber string `json:"latest_number" name:""`
	TimeStamp    string `json:"timestamp" name:""`
	BlockHash    string `json:"blockHash" name:""` // 块hash
	BlockNumber  string `json:"blockNumber" name:""`
	From         string `json:"from" name:""`
	//FromIsContract bool   `json:"from_is_contract" name:""`
	// Gas              string `json:"gas" name:""`
	GasPrice string `json:"gasPrice" name:""`
	Hash     string `json:"hash" name:""` // 交易 hash
	// Input string `json:"input" name:""`
	// Nonce string `json:"nonce" name:""`
	// R                string `json:"r" name:""`
	// S                string `json:"s" name:""`
	To string `json:"to" name:""`
	//ToIsContract bool   `json:"to_is_contract" name:""`
	// TransactionIndex string `json:"transactionIndex" name:""`
	// Type             string `json:"type" name:""`
	// V     string `json:"v" name:""`
	Value           string `json:"value" name:""`
	ContractAddress string `json:"contract_address" name:""` // 如果是 token transfer 就存在 ContractAddress
	IsContractToken bool   `json:"is_contract" name:""`      // 是否为合约交易
	FromBalance     string `json:"from_balance" name:""`     // 当前交易时用户余额
	ToBalance       string `json:"to_balance" name:""`
	TxType          string `json:"tx_type" name:""`    // 交易类型 外部交易，内部交易，代币交易
	EventType       string `json:"event_type" name:""` // 交易事件类型
}

func (sel *TransactionData) IsLatest() bool {
	return sel.LatestNumber == sel.BlockNumber
}

type TransactionRecord struct {
	ID              int64     `gorm:"primary_key"`
	CreatedAt       time.Time `json:"created_at" gorm:"column:created_at;COMMENT:" name:""`     // 记录创建时间
	BlockNumber     string    `json:"block_number" gorm:"column:block_number;COMMENT:" name:""` // 区块号
	BlockHash       string    `json:"block_hash" gorm:"column:block_hash;COMMENT:" name:""`
	TxHash          string    `json:"tx_hash" gorm:"column:tx_hash;index;varchar(255);COMMENT:" name:""`
	Timestamp       string    `json:"timestamp" gorm:"column:timestamp;COMMENT:" name:""`
	ContractAddress string    `json:"contract_address" gorm:"column:contract_address;COMMENT:" name:""`
	FromAddr        string    `json:"from_addr" gorm:"column:from_addr;COMMENT:" name:""`
	ToAddr          string    `json:"to_addr" gorm:"column:to_addr;COMMENT:" name:""`
	GasPrice        string    `json:"gas_price" gorm:"column:gas_price;COMMENT:" name:""`
	Value           string    `json:"value" gorm:"column:value;COMMENT:" name:""`
	FromAddrBalance string    `json:"balance" gorm:"column:balance;varchar(255);COMMENT:" name:""`
	ToAddrBalance   string    `json:"balance" gorm:"column:balance;varchar(255);COMMENT:" name:""`
}

type SplitTableInfo struct {
	ID           int64  `gorm:"primary_key"`
	TemplateName string `json:"template_name" gorm:"column:template_name;index;type:varchar(50);COMMENT:" name:""`
	TableName    string `json:"table_name" gorm:"column:table_name;COMMENT:" name:""`
}
