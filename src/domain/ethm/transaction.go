package ethm

import (
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"strconv"
)

type TransactionWriter struct {
	ethCli      ethrpc.EthRPC
	txRepo      repo.TransactionRepository
	withBalance bool
}

func NewTransactionWriter(ethCli ethrpc.EthRPC, txRepo repo.TransactionRepository) *TransactionWriter {
	txw := &TransactionWriter{
		ethCli: ethCli,
		txRepo: txRepo,
	}
	return txw
}

func (sel *TransactionWriter) SetWithBalance(withBalance bool) {
	sel.withBalance = withBalance
}

func (sel *TransactionWriter) TxWrite(txData *model.TransactionData) error {
	if sel.withBalance {
		blockNumber, _ := strconv.ParseInt(txData.BlockNumber, 10, 63)
		if txData.IsContractToken {
			txData.FromBalance, _ = sel.ethCli.GetContractBalanceByBlockNumber(txData.ContractAddress, txData.From, blockNumber)
			txData.ToBalance, _ = sel.ethCli.GetContractBalanceByBlockNumber(txData.ContractAddress, txData.To, blockNumber)
		} else {
			txData.FromBalance, _ = sel.ethCli.GetBalanceByBlockNumber(txData.From, blockNumber)
			txData.ToBalance, _ = sel.ethCli.GetBalanceByBlockNumber(txData.To, blockNumber)
		}
	}
	// 写入交易信息
	return sel.txRepo.CreateTransactionRecord(txData)
	//return nil
}
