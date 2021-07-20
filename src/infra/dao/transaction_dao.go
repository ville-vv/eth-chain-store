package dao

import (
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
)

// EthereumTransactionDao 以太坊交易
type EthereumTransactionDao struct {
}

func (sel *EthereumTransactionDao) CreateTransactionRecord(txData *model.TransactionData) error {
	fmt.Println("写入交易数据")
	return nil
}
