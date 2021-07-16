package model

const (
	TxTypeTransfer = "transfer"
)

type TransactionData struct {
	TimeStamp   string `json:"timestamp" name:""`
	BlockHash   string `json:"blockHash" name:""` // 块hash
	BlockNumber string `json:"blockNumber" name:""`
	From        string `json:"from" name:""`
	// Gas              string `json:"gas" name:""`
	// GasPrice         string `json:"gasPrice" name:""`
	Hash string `json:"hash" name:""` // 交易 hash
	// Input string `json:"input" name:""`
	// Nonce string `json:"nonce" name:""`
	// R                string `json:"r" name:""`
	// S                string `json:"s" name:""`
	To string `json:"to" name:""`
	// TransactionIndex string `json:"transactionIndex" name:""`
	// Type             string `json:"type" name:""`
	// V     string `json:"v" name:""`
	Value           string `json:"value" name:""`
	ContractAddress string `json:"contract_address" name:""`
	IsContract      bool   `json:"is_contract" name:""`
	Balance         string `json:"balance" name:""`
	TxType          string `json:"tx_type" name:""` // 交易类型
	IsErc20         bool   `json:"is_erc_20" name:""`
}
