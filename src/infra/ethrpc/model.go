package ethrpc

type RpcBlock struct {
	Hash            string            `json:"hash" name:""`
	Nonce           string            `json:"nonce" name:""`
	Number          string            `json:"number" name:""`
	ParentHash      string            `json:"parentHash" name:""`
	ReceiptsRoot    string            `json:"receiptsRoot" name:""`
	Size            string            `json:"size" name:""`
	TimeStamp       string            `json:"timestamp" name:""`
	TotalDifficulty string            `json:"totalDifficulty" name:""`
	Transactions    []*RpcTransaction `json:"transactions" name:""`
}

type RpcTransaction struct {
	BlockHash        string `json:"blockHash" name:""`
	BlockNumber      string `json:"blockNumber" name:""`
	From             string `json:"from" name:""`
	Gas              string `json:"gas" name:""`
	GasPrice         string `json:"gasPrice" name:""`
	Hash             string `json:"hash" name:""`
	Input            string `json:"input" name:""`
	Nonce            string `json:"nonce" name:""`
	R                string `json:"r" name:""`
	S                string `json:"s" name:""`
	To               string `json:"to" name:""`
	TransactionIndex string `json:"transactionIndex" name:""`
	Type             string `json:"type" name:""`
	V                string `json:"v" name:""`
	Value            string `json:"value" name:""`
}

type RpcTransactionReceiptLog struct {
	Address          string   `json:"address" name:""`
	BlockHash        string   `json:"blockHash" name:""`
	BlockNumber      string   `json:"blockNumber" name:""`
	Data             string   `json:"data" name:""`
	LogIndex         string   `json:"logIndex" name:""`
	Removed          string   `json:"removed" name:""`
	Topics           []string `json:"topics" name:""`
	TransactionHash  string   `json:"transactionHash" name:""`
	TransactionIndex string   `json:"transactionIndex" name:""`
}

type RpcTransactionReceipt struct {
	BlockHash         string                      `json:"blockHash" name:""`
	BlockNumber       string                      `json:"blockNumber" name:""`
	ContractAddress   string                      `json:"contract_address" name:""`
	CumulativeGasUsed string                      `json:"cumulativeGasUsed" name:""`
	From              string                      `json:"from" name:""`
	GasUsed           string                      `json:"gas_used" gorm:"column:gas_used;COMMENT:" name:""`
	Logs              []*RpcTransactionReceiptLog `json:"logs" name:""`
	LogsBloom         string                      `json:"logsBloom" name:""`
	Status            string                      `json:"status" name:""`
	To                string                      `json:"to" name:""`
	TransactionHash   string                      `json:"transactionHash" name:""`
	TransactionIndex  string                      `json:"transactionIndex" name:""`
	Type              string                      `json:"type" name:""`
}
