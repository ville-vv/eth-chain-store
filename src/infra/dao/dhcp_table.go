package dao

import (
	"errors"
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"gorm.io/gorm"
	"sync"
)

type dhcpTable struct {
	cntTable      string
	id            int64
	db            *gorm.DB
	lock          sync.Mutex
	maxNum        int64
	counter       int64  // 计数器
	templateTable string //
}

func newDhcpTable(db *gorm.DB, templateTable string) *dhcpTable {
	tb := &dhcpTable{
		db:            db,
		lock:          sync.Mutex{},
		maxNum:        20000000,
		counter:       0,
		templateTable: templateTable,
	}
	err := tb.intCntTxTable()
	if err != nil {
		panic(err)
	}
	if templateTable == "" {
		panic(errors.New("template table is empty"))
	}
	return tb
}

// intCntTxTable 初始化
func (sel *dhcpTable) intCntTxTable() error {
	db := sel.db
	tb := &model.SplitTableInfo{}
	err := db.Model(tb).Select("id, table_name").Where("template_name=?", sel.templateTable).Last(tb).Error
	if err != nil {
		if err.Error() != gorm.ErrRecordNotFound.Error() {
			return err
		}
		err = sel.createTxTable(sel.id)
		if err != nil {
			return err
		}
		sel.id = 1
	} else {
		sel.cntTable = tb.TableName
		sel.id = tb.ID
	}
	return nil
}

func (sel *dhcpTable) createTxTable(id int64) error {
	db := sel.db.Begin()
	defer db.Rollback()
	tbName := fmt.Sprintf("%s_%0.4d", sel.templateTable, id)
	err := db.Exec(fmt.Sprintf("create table if not exists %s like %s", tbName, sel.templateTable)).Error
	if err != nil {
		return err
	}
	tb := &model.SplitTableInfo{TemplateName: sel.templateTable, TableName: tbName}
	err = db.Create(tb).Error
	if err != nil {
		return err
	}
	sel.id = tb.ID
	sel.cntTable = tb.TableName
	return db.Commit().Error
}

// 表格覆盖
func (sel *dhcpTable) TbName() (string, error) {
	sel.lock.Lock()
	if sel.counter > sel.maxNum {
		err := sel.createTxTable(sel.id)
		if err != nil {
			return "", err
		}
		sel.id++
		sel.counter = 0
	}
	sel.lock.Unlock()
	return sel.cntTable, nil
}

func (sel *dhcpTable) Inc() {
	sel.lock.Lock()
	sel.counter++
	sel.lock.Unlock()
}

func (sel *dhcpTable) AllTable() ([]string, error) {
	db := sel.db.Model(&model.SplitTableInfo{})
	var tableList []model.SplitTableInfo
	err := db.Select("table_name").Where("template_name=?", sel.templateTable).Find(&tableList).Error
	if err != nil {
		return nil, err
	}

	tableNames := make([]string, 0, len(tableList))
	for _, val := range tableList {
		tableNames = append(tableNames, val.TableName)
	}
	return tableNames, nil
}
