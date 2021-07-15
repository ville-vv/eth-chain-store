package repo

import "github.com/ville-vv/eth-chain-store/src/infra/model"

type NormalAccountRepo struct {
}

func (sel *NormalAccountRepo) IsAccountExist(addr string) bool {
	return false
}

// UpdateContract
func (sel *NormalAccountRepo) UpdateContract(*model.AccountContent) error {
	return nil
}

// UpdateNative
func (sel *NormalAccountRepo) UpdateBalance(balance string) error {
	return nil
}

func (sel *NormalAccountRepo) CreateEthAccount() error {
	return nil
}

type ContractAccountRepo struct {
}

func (sel *ContractAccountRepo) IsAccountExist(addr string) bool {
	return false
}

// UpdateNative
func (sel *ContractAccountRepo) UpdateBalance(balance string) error {
	return nil
}
func (sel *ContractAccountRepo) CreateEthAccount() error {
	return nil
}
