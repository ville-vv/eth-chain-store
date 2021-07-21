package dao

import (
	"fmt"
	"github.com/ville-vv/eth-chain-store/src/common/conf"
	"github.com/ville-vv/vilgo/vstore"
	"gorm.io/gorm"
	"sync"
)

type DB interface {
	vstore.DB
	GetName() string
}

type MysqlDB struct {
	db   vstore.DB
	name string
}

func NewMysqlDB(db vstore.DB, name string) *MysqlDB {
	return &MysqlDB{db: db, name: name}
}

func (sel *MysqlDB) GetDB() *gorm.DB {
	return sel.db.GetDB()
}

func (sel *MysqlDB) GetName() string {
	return sel.name
}

func (sel *MysqlDB) ClearAllData() {
	sel.db.ClearAllData()
}

var dbPool *Erc20DbPool
var onceDo sync.Once

type DBRegister interface {
	RegisterDb() []string
}

func InitDb(rgt DBRegister) error {
	onceDo.Do(func() {
		dbPool = NewErc20DbPool()
		dbPool.RegisterDb(rgt.RegisterDb())
	})
	return nil
}

func GenDbName(symbol string) string {
	return fmt.Sprintf("eth_erc20_token_" + symbol)
}

func Erc20DB(key string) *gorm.DB {
	return dbPool.GetDB(key).GetDB()
}

// Erc 20 代币 数据库
type Erc20DbPool struct {
	mu  sync.Mutex
	dbs map[string]DB
}

func NewErc20DbPool() *Erc20DbPool {
	return &Erc20DbPool{mu: sync.Mutex{}, dbs: make(map[string]DB)}
}

func (sel *Erc20DbPool) RegisterDb(dbList []string) {
	for _, val := range dbList {
		dbCfg := conf.NewMysqlConf()
		dbCfg.DbName = val
		sel.add(val, NewMysqlDB(vstore.MakeDb(dbCfg), val))
	}
}

func (sel *Erc20DbPool) add(name string, db DB) {
	sel.mu.Lock()
	sel.dbs[name] = db
	sel.mu.Unlock()
}

func (sel *Erc20DbPool) GetDB(key string) DB {
	key = GenDbName(key)
	sel.mu.Lock()
	defer sel.mu.Unlock()
	db, ok := sel.dbs[key]
	if !ok {
		dbCfg := conf.NewMysqlConf()
		db = NewMysqlDB(vstore.MakeDb(dbCfg), key)
		if db != nil {
			sel.dbs[key] = db
		}
	}
	return db
}
