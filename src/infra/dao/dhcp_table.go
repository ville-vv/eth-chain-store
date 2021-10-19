package dao

import (
	"errors"
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/infra/model"
	"github.com/ville-vv/vilgo/vlog"
	"gorm.io/gorm"
	"strconv"
	"strings"
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
	if templateTable == "" {
		panic(errors.New("template table is empty"))
	}
	tb := &dhcpTable{
		db:            db,
		lock:          sync.Mutex{},
		maxNum:        10000000,
		counter:       0,
		templateTable: templateTable,
	}

	return tb
}

func (sel *dhcpTable) Init() error {
	return sel.intCntTxTable()
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
		strList := strings.Split(tb.TableName, "_")
		if len(strList)> 0{
			tbSeqStr := strList[len(strList)-1]
			tbSeq, _ := strconv.Atoi(tbSeqStr)
			sel.id = int64(tbSeq)
		}else{
			sel.id = tb.ID
		}
		if err := sel.count(); err != nil {
			return err
		}
	}

	return nil
}

// count 统计当前表的总数
func (sel *dhcpTable) count() error {
	db := sel.db
	var count int64
	err := db.Table(sel.cntTable).Count(&count).Error
	if err != nil {
		return err
	}
	sel.counter = count
	vlog.INFO("current table %s total record %d", sel.cntTable, count)
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

// TbName 表格覆盖
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
