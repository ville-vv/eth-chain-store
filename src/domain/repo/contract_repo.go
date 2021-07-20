package repo

import (
	"github.com/ville-vv/eth-chain-store/src/infra/dao"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/vilgo/vlog"
)

// ContractRepo 合约信息处理
type ContractRepo struct {
	contractDao *dao.EthereumDao
}

func (sel *ContractRepo) IsContractExist(addr string) bool {
	var contractInfo model.TbContractAddress
	if err := sel.contractDao.QueryContractInfo(addr, &contractInfo); err != nil {
		vlog.ERROR("")
		return true
	}
	if contractInfo.ID > 0 {
		return true
	}
	return false
}

func (sel *ContractRepo) CreateContract(ct *model.ContractContent) error {
	return sel.contractDao.CreateContractRecord(ct)
}
