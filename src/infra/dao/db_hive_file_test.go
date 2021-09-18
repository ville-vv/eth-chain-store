package dao

import (
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/common/utils"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"testing"
	"time"
)

func TestWalkField(t *testing.T) {
	record := &model.TransactionRecord{
		ID:              0,
		CreatedAt:       time.Time{},
		BlockNumber:     "3000003",
		BlockHash:       "0x18bdeb8f6ffe881a1320fc5dfd1711f618c77b8381805a5d78e148c237a33497",
		TxHash:          "0xee45955b53628c1492dce697607ec26122099f1f5b1b52e92fa62eb805321b5b",
		TxTime:          "2017-01-15 18:11:36",
		ContractAddress: "1",
		FromAddr:        "0xa8f769b88d6d74fb2bd3912f6793f75625228baf",
		ToAddr:          "0x4eb36b2a7e3c41b3970cfe179019b9b4093890f1",
		GasPrice:        "000",
		Value:           "000",
		FromAddrBalance: "000",
		ToAddrBalance:   "000",
	}
	fmt.Println(record.String())

	utils.WalkField(nil, func(name string, val interface{}) {
		fmt.Println(name, val)
	})
}
