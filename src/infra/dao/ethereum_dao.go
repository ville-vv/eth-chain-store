package dao

import (
	"github.com/pkg/errors"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"gorm.io/gorm"
)

type EthereumDao struct {
	db               DB
	contractRecordTb *dhcpTable
	contractBindTb   *dhcpTable
	normalAccountTb  *dhcpTable
}

func NewEthereumDao(db DB) *EthereumDao {
	return &EthereumDao{
		db:               db,
		contractRecordTb: newDhcpTable(db.GetDB(), "contract_address_records"),
		contractBindTb:   newDhcpTable(db.GetDB(), "contract_account_binds"),
		normalAccountTb:  newDhcpTable(db.GetDB(), "ethereum_accounts"),
	}
}

// QueryContractInfo 查询合约信息
func (sel *EthereumDao) QueryContractInfo(addr string, contractInfo *model.ContractAddressRecord) error {
	tbNames, err := sel.contractRecordTb.AllTable()
	if err != nil {
		return err
	}
	var db = sel.db.GetDB().Where("address=?", addr)
	for _, tbName := range tbNames {
		err = db.Table(tbName).First(contractInfo).Error
		if err != nil {
			if err.Error() != gorm.ErrRecordNotFound.Error() {
				return err
			}
		}
		if contractInfo.ID > 0 {
			break
		}
	}
	return nil
}

// CreateContractRecord 创建合约记录
func (sel *EthereumDao) CreateContractRecord(contractCtx *model.ContractContent) error {
	tbName, err := sel.contractRecordTb.TbName()
	if err != nil {
		return errors.Wrap(err, "ethereum dao create contract record")
	}
	if err = sel.db.GetDB().Table(tbName).Create(&model.ContractAddressRecord{
		ContractContent: *contractCtx,
	}).Error; err != nil {
		return err
	}
	sel.contractRecordTb.Inc()
	return nil
}

// QueryBindContractAccount 查询合约绑定账户信息
func (sel *EthereumDao) QueryBindContractAccount(addr, contractAddr string, bindInfo *model.ContractAccountBind) error {
	return nil
}

// CreateContractAccount 普通钱包地址与合约地址绑定
func (sel *EthereumDao) CreateContractAccount(contractAccount *model.ContractAccountBind) error {
	//fmt.Println("绑定钱包地址与合约地址", contractAccount)
	tbName, err := sel.contractBindTb.TbName()
	if err != nil {
		return errors.Wrap(err, "ethereum dao create contract account bind")
	}
	if err = sel.db.GetDB().Table(tbName).Create(contractAccount).Error; err != nil {
		return err
	}
	sel.contractBindTb.Inc()
	return nil
}

func (sel *EthereumDao) UpdateContractAccountBalance(addr string, contractAddr string, balance string) (bool, error) {
	tables, err := sel.contractBindTb.AllTable()
	if err != nil {
		return false, err
	}
	var db = sel.db.GetDB().Where("address=? and contract_address=?", addr, contractAddr)
	var contractAccount model.ContractAccountBind
	for _, tbName := range tables {
		if err = db.Table(tbName).Select("id, address, balance").First(&contractAccount).Error; err != nil {
			if err.Error() != gorm.ErrRecordNotFound.Error() {
				return false, errors.Wrap(err, "ethereum dao for update and query contract account exist")
			}
		}
		if contractAccount.ID > 0 && contractAccount.Balance != balance {
			return true, db.Table(tbName).Update("balance", balance).Error
		}
	}
	return false, nil
}

//========================================================================================================
// QueryNormalAccount 查询普通的以太坊账户地址信息
func (sel *EthereumDao) QueryNormalAccount(addr string, info *model.EthereumAccount) error {
	return nil
}

// UpdateNativeBalance 更新以太坊
func (sel *EthereumDao) UpdateNormalAccountBalance(addr string, balance string) (bool, error) {
	tables, err := sel.normalAccountTb.AllTable()
	if err != nil {
		return false, err
	}
	var db = sel.db.GetDB().Select("id, address,balance").Where("address=?", addr)
	var ethereumAccount model.EthereumAccount
	for _, tbName := range tables {
		if err = db.Table(tbName).First(&ethereumAccount).Error; err != nil {
			if err.Error() != gorm.ErrRecordNotFound.Error() {
				return false, errors.Wrap(err, "ethereum dao for update and query normal account exist")
			}
		}
		if ethereumAccount.ID > 0 && ethereumAccount.Balance != balance {
			return true, sel.db.GetDB().Table(tbName).Where("address=?", addr).Update("balance", balance).Error
		}
	}
	return false, nil
}

// CreateNormalAccount 创建以太坊账户
func (sel *EthereumDao) CreateNormalAccount(normalAccount *model.EthereumAccount) error {
	tbName, err := sel.normalAccountTb.TbName()
	if err != nil {
		return errors.Wrap(err, "ethereum dao create ethereum account")
	}
	if err = sel.db.GetDB().Table(tbName).Create(normalAccount).Error; err != nil {
		return err
	}
	sel.normalAccountTb.Inc()
	return nil
}
