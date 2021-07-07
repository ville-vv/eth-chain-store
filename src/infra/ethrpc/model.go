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
