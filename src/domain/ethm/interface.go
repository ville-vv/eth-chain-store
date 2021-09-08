package ethm

import "github.com/ville-vv/eth-chain-store/src/infra/model"

type TxWriter interface {
	TxWrite(txData *model.TransactionData) error
}

type TxWriteFun func(*model.TransactionData) error

func (tw TxWriteFun) TxWrite(txData *model.TransactionData) error {
	return tw(txData)
}
