package dao

import (
	"context"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/vilgo/vlog"
	"os"
)

type TransactionHiveDataFile struct {
	normalTxFs   *os.File
	contractTxFs *os.File
}

func (sel *TransactionHiveDataFile) Scheme() string {
	return "TransactionHiveDataFile"
}

func (sel *TransactionHiveDataFile) Init() error {
	return nil
}

func (sel *TransactionHiveDataFile) Start() error {
	return nil
}

func (sel *TransactionHiveDataFile) Exit(ctx context.Context) error {
	_ = sel.contractTxFs.Close()
	_ = sel.normalTxFs.Close()
	return nil
}

func NewHiveDataFile(normalTxFile string, contractTxFile string) *TransactionHiveDataFile {
	normalTxFs, err := os.OpenFile(normalTxFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err.Error() + normalTxFile)
	}

	contractTxFs, err := os.OpenFile(contractTxFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err.Error() + contractTxFile)
	}

	return &TransactionHiveDataFile{
		normalTxFs:   normalTxFs,
		contractTxFs: contractTxFs,
	}
}

func (sel *TransactionHiveDataFile) CreateTransactionRecord(txData *model.TransactionData) error {
	vlog.INFO("创建交易数据")
	if txData.IsContractToken {
		return sel.InsertContractData(txData)
	}
	return sel.InsertNormalData(txData)
}

func (sel *TransactionHiveDataFile) InsertNormalData(txData *model.TransactionData) error {
	return nil
}

func (sel *TransactionHiveDataFile) InsertContractData(txData *model.TransactionData) error {
	return nil
}
