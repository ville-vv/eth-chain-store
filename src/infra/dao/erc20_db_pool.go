package dao

import (
	"github.com/ville-vv/vilgo/vstore"
	"gorm.io/gorm"
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
