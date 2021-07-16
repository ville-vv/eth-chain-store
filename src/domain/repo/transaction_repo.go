package repo

import "github.com/ville-vv/eth-chain-store/src/infra/model"

type TransactionRepo struct {
}

func (sel *TransactionRepo) CreateTransactionRecord(txData *model.TransactionData) error {
	return nil
}
