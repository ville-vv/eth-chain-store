package dao

import (
	"context"
	"github.com/ville-vv/vilgo/vfile"
	"github.com/ville-vv/vilgo/vlog"
	"os"
	"path"
	"sync"
	"time"
)

type cacheElm struct {
	TableName string
	Record    interface{}
}

type cacheList []*cacheElm

type DbCache struct {
	sync.RWMutex
	db         DB
	cachePool  [3]cacheList
	elmList    cacheList
	poolIdx    int
	isStop     bool
	stopCh     chan int
	wrInterval int
	perFile    *os.File
}

func NewDbCache(perFile string, wrInterval int, db DB) *DbCache {
	if wrInterval <= 0 {
		wrInterval = 1
	}

	cachePool := [3]cacheList{make(cacheList, 0, 100000), make(cacheList, 0, 100000), make(cacheList, 0, 100000)}
	dirPath := path.Dir(perFile)
	if !vfile.PathExists(dirPath) {
		err := os.Mkdir(dirPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	f, err := os.OpenFile(perFile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	d := &DbCache{
		RWMutex:    sync.RWMutex{},
		db:         db,
		cachePool:  cachePool,
		elmList:    cachePool[0],
		isStop:     false,
		stopCh:     make(chan int),
		wrInterval: wrInterval,
		perFile:    f,
	}

	return d
}

func (sel *DbCache) Scheme() string {
	return "DbCache"
}

func (sel *DbCache) Init() error {
	return nil
}

func (sel *DbCache) Start() error {
	go sel.loopInsert()
	vlog.INFO("数据写入缓存启动")
	return nil
}

func (sel *DbCache) Exit(ctx context.Context) error {
	close(sel.stopCh)
	sel.perFile.Close()
	time.Sleep(time.Second)
	return nil
}

func (sel *DbCache) loopInsert() {
	tmr := time.NewTicker(time.Second * time.Duration(sel.wrInterval))
	for {
		select {
		case <-tmr.C:
			sel.saveCacheToDb()
		case <-sel.stopCh:
			sel.saveCacheToDb()
			return
		}
	}
}

func (sel *DbCache) saveCacheToDb() {
	sel.Lock()
	waitSaveList := sel.cachePool[sel.poolIdx]
	sel.cachePool[sel.poolIdx] = sel.cachePool[sel.poolIdx][:0]
	sel.poolIdx++
	if sel.poolIdx >= 3 {
		sel.poolIdx = 0
	}
	sel.Unlock()

	tempCacheMap := make(map[string][]interface{})
	for _, v := range waitSaveList {
		if v == nil {
			break
		}
		lst, ok := tempCacheMap[v.TableName]
		if !ok {
			lst = make([]interface{}, 0, len(waitSaveList))
		}
		lst = append(lst, v.Record)
		tempCacheMap[v.TableName] = lst
	}
	for tbName, v := range tempCacheMap {
		//vlog.INFO("插入到数据库 %s %d", tbName, len(v))
		db := sel.db.GetDB().Begin()
		sqlStr := BatchInsertToSqlStr(tbName, v)
		err := db.Exec(sqlStr).Error
		//err := DoBatchInsert(tbName, v, sel.db.GetDB())
		if err != nil {
			vlog.ERROR("save data to db table %s len:%d error %s", tbName, len(v), err.Error())
			_, _ = sel.perFile.WriteString(sqlStr + ";\n")
			db.Rollback()
			continue
		}
		db.Commit()
	}
	tempCacheMap = nil
}

func (sel *DbCache) Insert(tableName string, val interface{}) error {
	sel.Lock()
	sel.cachePool[sel.poolIdx] = append(sel.cachePool[sel.poolIdx], &cacheElm{TableName: tableName, Record: val})
	sel.Unlock()
	return nil
}

func (sel *DbCache) Select(fn func(tbName string, val interface{})) error {
	sel.RLock()
	for _, val := range sel.elmList {
		fn(val.TableName, val.Record)
	}
	sel.RUnlock()
	return nil
}