package dao

import (
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"sync"
)

type Erc20TokenConfigDao struct {
	st *sync.Map
	db DB
}

func NewErc20TokenConfigDao(db DB) *Erc20TokenConfigDao {
	e := &Erc20TokenConfigDao{db: db}
	err := e.Loading()
	if err != nil {
		panic(err)
	}
	return e
}

func (sel *Erc20TokenConfigDao) TokenIsValid(addr string) bool {
	_, ok := sel.st.Load(addr)
	if !ok {
		// 如果没有就去数据库中查询
		var e2tc model.Erc20TokenConfigContent
		db := sel.db.GetDB().Table("eth_contract_erc20_config").Select("symbol")
		err := db.Where("address=?", addr).First(&e2tc).Error
		if err != nil {
			return false
		}
		e2tc.Address = addr
		sel.st.Store(addr, &e2tc)
	}
	return true
}

func (sel *Erc20TokenConfigDao) Loading() error {
	// 如果没有就去数据库中查询
	var e2tcList []*model.Erc20TokenConfigContent
	db := sel.db.GetDB().Table("eth_contract_erc20_config").Select("address,symbol")
	err := db.Find(&e2tcList).Error
	if err != nil {
		return err
	}
	for i, val := range e2tcList {
		sel.st.Store(val.Address, e2tcList[i])
	}
	return nil
}

func (sel *Erc20TokenConfigDao) RegisterDb() []string {
	var dbs []string
	sel.st.Range(func(key, value interface{}) bool {
		erc20Cfg, Ok := value.(*model.Erc20TokenConfigContent)
		if Ok {
			dbs = append(dbs, GenDbName(erc20Cfg.Symbol))
		}
		return true
	})
	return dbs
}

func (sel *Erc20TokenConfigDao) Range(rf func(model.Erc20TokenConfigContent)) {
	sel.st.Range(func(key, value interface{}) bool {
		erc20Cfg, Ok := value.(*model.Erc20TokenConfigContent)
		if Ok {
			rf(*erc20Cfg)
		}
		return true
	})
	return
}
