package ethm

import "github.com/ville-vv/eth-chain-store/src/infra/model"

type TxWriter interface {
	TxWrite(txData *model.TransactionData) error
}

type TxWriterFilter interface {
	Filter(addr string) (err error)
}
