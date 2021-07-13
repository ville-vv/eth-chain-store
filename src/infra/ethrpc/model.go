package ethrpc

import "github.com/ville-vv/eth-chain-store/src/common/go-ethereum/common"

type RpcBlockHeader struct {
	Hash            string `json:"hash" name:""`
	Nonce           string `json:"nonce" name:""`
	Number          string `json:"number" name:""`
	ParentHash      string `json:"parentHash" name:""`
	ReceiptsRoot    string `json:"receiptsRoot" name:""`
	Size            string `json:"size" name:""`
	TimeStamp       string `json:"timestamp" name:""`
	TotalDifficulty string `json:"totalDifficulty" name:""`
}

type RpcBlock struct {
	RpcBlockHeader
	Transactions []*RpcTransaction `json:"transactions" name:""`
}

type RpcTransaction struct {
	BlockHash        string `json:"blockHash" name:""`
	BlockNumber      string `json:"blockNumber" name:""`
	From             string `json:"from" name:""`
	Gas              string `json:"gas" name:""`
	GasPrice         string `json:"gasPrice" name:""`
	Hash             string `json:"hash" name:""`
	Input            Input  `json:"input" name:""`
	Nonce            string `json:"nonce" name:""`
	R                string `json:"r" name:""`
	S                string `json:"s" name:""`
	To               string `json:"to" name:""`
	TransactionIndex string `json:"transactionIndex" name:""`
	Type             string `json:"type" name:""`
	V                string `json:"v" name:""`
	Value            string `json:"value" name:""`
}

func (t *RpcTransaction) IsContract() bool {
	return t.Input.IsContract()
}

func (t *RpcTransaction) IsTransfer() bool {
	return t.Input.IsTransfer()
}

type RpcTransactionReceiptLog struct {
	Address          string   `json:"address" name:""`
	BlockHash        string   `json:"blockHash" name:""`
	BlockNumber      string   `json:"blockNumber" name:""`
	Data             string   `json:"data" name:""`
	LogIndex         string   `json:"logIndex" name:""`
	Removed          bool     `json:"removed" name:""`
	Topics           []string `json:"topics" name:""`
	TransactionHash  string   `json:"transactionHash" name:""`
	TransactionIndex string   `json:"transactionIndex" name:""`
}

func (sel *RpcTransactionReceiptLog) IsTransfer() bool {
	if len(sel.Topics) == 3 {
		if GetNewLabelFromSignature(sel.Topics[0]) == Transfer {
			return true
		}
	}
	return false
}

func (sel *RpcTransactionReceiptLog) From() string {
	if len(sel.Topics) == 3 {
		return common.HexToAddress(sel.Topics[1]).String()
	}
	return ""
}

func (sel *RpcTransactionReceiptLog) To() string {
	if len(sel.Topics) == 3 {
		return common.HexToAddress(sel.Topics[2]).String()
	}
	return ""
}

func (sel *RpcTransactionReceiptLog) Value() string {
	return common.HexToHash(sel.Data).Big().String()
}

type RpcTransactionReceipt struct {
	BlockHash         string                      `json:"blockHash" name:""`
	BlockNumber       string                      `json:"blockNumber" name:""`
	ContractAddress   string                      `json:"contractAddress" name:""`
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