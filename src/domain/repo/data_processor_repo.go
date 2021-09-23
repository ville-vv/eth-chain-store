package repo

import "github.com/ville-vv/eth-chain-store/src/infra/model"

type DataProcessorRepository interface {
}

type DataProcessorRepositoryImpl struct {
}

func (sel *DataProcessorRepositoryImpl) GetTxData() ([]*model.TransactionRecord, error) {
	return nil, nil
}

func (sel *DataProcessorRepositoryImpl) Finish() error {
	return nil
}
