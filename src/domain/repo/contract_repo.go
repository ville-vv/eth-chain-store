package repo

import (
	"github.com/ville-vv/eth-chain-store/src/infra/dao"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/vilgo/vlog"
)

type ContractRepository interface {
	CreateContract(ct *model.ContractContent) error
	IsContractExist(addr string) bool
}

// ContractRepo 合约信息处理
type ContractRepo struct {
	contractDao *dao.EthereumDao
}

func NewContractRepo(contractDao *dao.EthereumDao) *ContractRepo {
	return &ContractRepo{contractDao: contractDao}
}

func (sel *ContractRepo) IsContractExist(addr string) bool {
	var contractInfo model.ContractAddressRecord
	if err := sel.contractDao.QueryContractInfo(addr, &contractInfo); err != nil {
		vlog.ERROR("ContractRepo.IsContractExist query contract information failed %s %s", addr, err.Error())
		return true
	}
	if contractInfo.ID > 0 {
		return true
	}
	return false
}

func (sel *ContractRepo) CreateContract(ct *model.ContractContent) error {
	// vlog.DEBUG("create contract information Address:%s Symbol:%s TotalSupply:%s", ct.Address, ct.Symbol, ct.TotalSupply)
	return sel.contractDao.CreateContractRecord(ct)
}
