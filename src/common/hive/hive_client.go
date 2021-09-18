package hive

//
import (
	"context"
	"fmt"
	"github.com/beltran/gohive"
	"github.com/pkg/errors"
	"reflect"
)

type HiveConfigOption interface {
	GetHost() string
	GetPort() int
	GetDBName() string
	GetAuthMode() string
	GetUserName() string
	GetPassword() string
}

type HiveCLI struct {
	opt        HiveConfigOption
	connCfg    gohive.ConnectConfiguration
	conn       *gohive.Connection
	defaultCtx context.Context
}

func New(opt HiveConfigOption) (*HiveCLI, error) {
	cfg := gohive.NewConnectConfiguration()
	cfg.Service = "hive"
	cfg.Database = opt.GetDBName()

	ghConn, err := gohive.Connect(opt.GetHost(), opt.GetPort(), opt.GetAuthMode(), cfg)
	if err != nil {
		return nil, err
	}
	hcli := &HiveCLI{
		conn:       ghConn,
		opt:        opt,
		connCfg:    gohive.ConnectConfiguration{},
		defaultCtx: context.Background(),
	}
	return hcli, err
}

func (sel *HiveCLI) Close() {
	sel.conn.Close()
}

func (sel *HiveCLI) Find(stm string, res interface{}) error {
	if res == nil {
		return errors.New("res is nil")
	}

	objType := reflect.TypeOf(res)
	objValPtr := reflect.ValueOf(res)

	if objType.Kind() != reflect.Ptr {
		return errors.New("res is must ptr")
	}
	objVal := objValPtr.Elem()
	objType = objType.Elem()

	var dataMap map[string]interface{}
	//dataMap = map[string]interface{}{
	//	"top1000_erc20_token.index":                  int32(1234),
	//	"top1000_erc20_token.token_contract_address": "0x0000000000004946c0e9F43F4Dee607b0eF1fA1c"}
	cursor := sel.conn.Cursor()
	defer cursor.Close()
	cursor.Exec(sel.defaultCtx, stm)
	if cursor.Err != nil {
		return errors.Wrap(cursor.Err, "hive exec")
	}
	for cursor.HasMore(sel.defaultCtx) {
		if cursor.Err != nil {
			return errors.Wrap(cursor.Err, "hive exec has more")
		}
		dataMap = cursor.RowMap(sel.defaultCtx)
		if cursor.Err != nil {
			return errors.Wrap(cursor.Err, "hive exec row map")
		}

		if objType.Kind() == reflect.Slice {
			elm := objType.Elem()
			if elm.Kind() == reflect.Ptr {
				elm = elm.Elem()
			}
			tempObj := mapToStru(elm, dataMap)
			// 给外部变量赋值
			objVal = reflect.Append(objVal, tempObj)

		} else if objType.Kind() == reflect.Struct {
			tempObj := mapToStru(objType, dataMap)
			objVal.Set(tempObj)
			return nil
		}

	}
	objValPtr.Elem().Set(objVal)
	return nil
}

func (sel *HiveCLI) Count(tableName string) int64 {
	cursor := sel.conn.Cursor()
	defer cursor.Close()

	cursor.Exec(sel.defaultCtx, fmt.Sprintf("select count(*) as total from %s", tableName))

	datas := cursor.RowMap(sel.defaultCtx)

	total, _ := datas["total"].(int64)

	return total
}

func (sel *HiveCLI) Exec(stm string) error {
	cursor := sel.conn.Cursor()
	defer cursor.Close()
	cursor.Exec(sel.defaultCtx, stm)
	if cursor.Err != nil {
		return errors.Wrap(cursor.Err, "hive exec")
	}
	return nil
}

func mapToStru(oType reflect.Type, dataMap map[string]interface{}) reflect.Value {
	var columnName string
	// 通过反射创建一个对象
	tempObjPtr := reflect.New(oType)
	tempObj := tempObjPtr.Elem()
	// 给该对象赋值
	for i := 0; i < tempObj.Type().NumField(); i++ {
		if tempObj.Type().Field(i).Tag.Get("sql") == "-" {
			continue
		}
		columnName = tempObj.Type().Field(i).Tag.Get("column")
		if columnName == "" {
			columnName = ToDbName(tempObj.Type().Name())
		}
		tempData := dataMap[columnName]
		if tempData == nil {
			continue
		}
		tempDataVal := reflect.ValueOf(tempData)

		amd := tempObj.Field(i)
		amd.Set(tempDataVal)
	}
	return tempObjPtr
}
