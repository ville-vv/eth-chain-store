package ethrpc

import (
	"github.com/ville-vv/eth-chain-store/src/common/go-eth/common"
	"strconv"
	"time"
)

type EthBlockHeader struct {
	LatestBlockNumber string `json:"latest_block_number" name:""`
	Hash              string `json:"hash" name:""`
	Nonce             string `json:"nonce" name:""`
	Number            string `json:"number" name:""`
	ParentHash        string `json:"parentHash" name:""`
	ReceiptsRoot      string `json:"receiptsRoot" name:""`
	Size              string `json:"size" name:""`
	TimeStamp         string `json:"timestamp" name:""`
	TotalDifficulty   string `json:"totalDifficulty" name:""`
}

func (sel *EthBlockHeader) TimeStampFormatTmString() string {
	sec, _ := strconv.ParseInt(common.HexToHash(sel.TimeStamp).Big().String(), 10, 64)
	return time.Unix(sec, 0).Format("2006-01-02 15:04:05")
}

type EthBlock struct {
	EthBlockHeader
	Transactions []*EthTransaction `json:"transactions" name:""`
}

type EthTransaction struct {
	BlockHash   string `json:"blockHash" name:""` // 块hash
	BlockNumber string `json:"blockNumber" name:""`
	From        string `json:"from" name:""`
	// Gas              string `json:"gas" name:""`
	GasPrice string `json:"gasPrice" name:""`
	Hash     string `json:"hash" name:""` // 交易hash
	Input    Input  `json:"input" name:""`
	// Nonce            string `json:"nonce" name:""`
	// R                string `json:"r" name:""`
	// S                string `json:"s" name:""`
	To string `json:"to" name:""`
	// TransactionIndex string `json:"transactionIndex" name:""`
	// Type             string `json:"type" name:""`
	// V                string `json:"v" name:""`
	Value string `json:"value" name:""`
}

func (t *EthTransaction) BlockNumberToBig() string {
	return common.HexToHash(t.BlockNumber).Big().String()
}

func (t *EthTransaction) GasPriceToBig() string {
	return common.HexToHash(t.GasPrice).Big().String()
}

func (t *EthTransaction) IsContractToken() bool {
	return t.Input.IsContractToken()
}

// IsTransfer 判断是否为Erc20协议的合约交易转账
func (t *EthTransaction) IsTransfer() bool {
	return t.Input.IsErc20Transfer()
}

type EthTransactionReceiptLog struct {
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

func (sel *EthTransactionReceiptLog) BlockNumberToBig() string {
	return common.HexToHash(sel.BlockNumber).Big().String()
}

func (sel *EthTransactionReceiptLog) IsTransfer() bool {
	if len(sel.Topics) == 3 {
		if GetNewLabelFromSignature(sel.Topics[0]) == Transfer {
			return true
		}
	}
	return false
}

func (sel *EthTransactionReceiptLog) From() string {
	if len(sel.Topics) == 3 {
		return common.HexToAddress(sel.Topics[1]).String()
	}
	return ""
}

func (sel *EthTransactionReceiptLog) To() string {
	if len(sel.Topics) == 3 {
		return common.HexToAddress(sel.Topics[2]).String()
	}
	return ""
}

func (sel *EthTransactionReceiptLog) Value() string {
	return common.HexToHash(sel.Data).Big().String()
}

type EthTransactionReceipt struct {
	BlockHash         string                      `json:"blockHash" name:""`
	BlockNumber       string                      `json:"blockNumber" name:""`
	ContractAddress   string                      `json:"contractAddress" name:""`
	CumulativeGasUsed string                      `json:"cumulativeGasUsed" name:""`
	From              string                      `json:"from" name:""`
	GasUsed           string                      `json:"gas_used" gorm:"column:gas_used;COMMENT:" name:""`
	Logs              []*EthTransactionReceiptLog `json:"logs" name:""`
	LogsBloom         string                      `json:"logsBloom" name:""`
	Status            string                      `json:"status" name:""`
	To                string                      `json:"to" name:""`
	TransactionHash   string                      `json:"transactionHash" name:""`
	TransactionIndex  string                      `json:"transactionIndex" name:""`
	Type              string                      `json:"type" name:""`
}
