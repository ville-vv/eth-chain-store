package repo

import "github.com/ville-vv/eth-chain-store/src/infra/model"

type ContractRepo struct {
}

func (sel *ContractRepo) IsContractExist(addr string) bool {
	return false
}

func (sel *ContractRepo) UpdateContract(ct *model.ContractAccount) error {
	return nil
}

func (sel *ContractRepo) CreateContract(ct *model.ContractAccount) error {
	return nil
}
