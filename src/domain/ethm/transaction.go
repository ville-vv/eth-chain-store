package ethm

import (
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/vilgo/vlog"
)

type TransactionWriter struct {
	ethCli ethrpc.EthRPC
	txRepo *repo.TransactionRepo
}

func NewTransactionWriter(ethCli ethrpc.EthRPC, txRepo *repo.TransactionRepo) *TransactionWriter {
	txw := &TransactionWriter{
		ethCli: ethCli,
		txRepo: txRepo,
	}
	return txw
}

func (sel *TransactionWriter) TxWrite(txData *model.TransactionData) error {
	vlog.DEBUG("tx writer to transaction %s", txData.Hash)
	// 写入交易信息
	return sel.txRepo.CreateTransactionRecord(txData)
}
