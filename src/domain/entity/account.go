package entity

import (
	"github.com/ville-vv/eth-chain-store/src/domain/repo"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
)

type EthAccount struct {
	model.EthereumAccount
}

type ContractAccount struct {
	accountRepo repo.ContractRepository
	model.EthereumAccount
	ContractAddress string `json:"contract_address" gorm:"column:contract_address;type:varchar(128);index:idx_cab_addr_crtaddr;COMMENT:" name:""` // 合约代币地址
	Symbol          string `json:"symbol" gorm:"column:symbol;type:varchar(128);COMMENT:" name:""`                                                // 合约代币标识
}

func (sel *ContractAccount) CreateRecord() error {
	return nil
}
