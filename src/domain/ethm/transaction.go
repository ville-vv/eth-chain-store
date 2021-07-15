package ethm

import (
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/ethrpc"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
)

type TransactionWriter struct {
	ethCli  ethrpc.EthRPC
	actRepo repo.NormalAccountRepo
}

func (sel *TransactionWriter) TxWrite(txData *model.TransactionData) error {
	// 查询该地址是否存在，写入地址
	return nil
}
